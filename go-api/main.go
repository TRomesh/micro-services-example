package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
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
			"message": "pong 🏓",
		})
	})

	// Start the server
	r.Run(fmt.Sprintf(":%s", PORT))
}
