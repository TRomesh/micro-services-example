package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var mongoClient *mongo.Client
var mongoCollection *mongo.Collection

func initMongoDB() {
	// Read MongoDB connection details from environment variables
	MONGO_URI := viper.GetString("MONGO_URI")
	DATABASE_NAME := viper.GetString("DATABASE_NAME")
	COLLECTION_NAME := viper.GetString("COLLECTION_NAME")

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	failOnError(err, "Failed to connect to MongoDB")
	mongoClient = client

	// Select database and collection
	mongoCollection = client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	log.Printf("Connected to MongoDB database: %s, collection: %s", DATABASE_NAME, COLLECTION_NAME)
}

func createDocument(c *gin.Context) {
	var document bson.M
	if err := c.BindJSON(&document); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	result, err := mongoCollection.InsertOne(context.Background(), document)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create document"})
		return
	}
	c.JSON(201, gin.H{"insertedID": result.InsertedID})
}

func readDocument(c *gin.Context) {
	id := c.Param("id")
	var document bson.M
	err := mongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&document)
	if err != nil {
		c.JSON(404, gin.H{"error": "Document not found"})
		return
	}
	c.JSON(200, document)
}

func updateDocument(c *gin.Context) {
	id := c.Param("id")
	var update bson.M
	if err := c.BindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	result, err := mongoCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update document"})
		return
	}
	c.JSON(200, gin.H{"matchedCount": result.MatchedCount, "modifiedCount": result.ModifiedCount})
}

func deleteDocument(c *gin.Context) {
	id := c.Param("id")
	result, err := mongoCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete document"})
		return
	}
	c.JSON(200, gin.H{"deletedCount": result.DeletedCount})
}

func main() {
	r := gin.Default()

	// Configure Viper to read environment variables
	viper.AutomaticEnv()

	// Read environment variables
	PORT := viper.GetString("APP_PORT")
	QUEUE_NAME := viper.GetString("QUEUE_NAME")
	QUEUE_USER := viper.GetString("QUEUE_USER")
	QUEUE_HOST := viper.GetString("QUEUE_HOST")
	QUEUE_PASSWORD := viper.GetString("QUEUE_PASSWORD")

	// Initialize MongoDB connection
	initMongoDB()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %s", err)
		}
	}()

	// Log connection string (omit sensitive data in production)
	queueStringConnection := fmt.Sprintf("amqp://%s:%s@%s", QUEUE_USER, QUEUE_PASSWORD, QUEUE_HOST)
	fmt.Println("Connecting to RabbitMQ with: " + queueStringConnection)

	// Establish RabbitMQ connection
	conn, err := amqp.Dial(queueStringConnection)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open RabbitMQ channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare RabbitMQ queue
	q, err := ch.QueueDeclare(
		QUEUE_NAME, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Process messages in a goroutine
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	// Define HTTP routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "go-api pong 🏓",
		})
	})
	r.POST("/documents", createDocument)
	r.GET("/documents/:id", readDocument)
	r.PUT("/documents/:id", updateDocument)
	r.DELETE("/documents/:id", deleteDocument)

	// Start the server
	r.Run(fmt.Sprintf(":%s", PORT))
}
