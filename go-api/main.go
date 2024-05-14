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
	viper.SetConfigFile(".env")
    viper.ReadInConfig()
	PORT := viper.Get("APP_PORT")
	USER_QUEUE := viper.GetString("USER_QUEUE")
	QUEUE_USER := viper.GetString("QUEUE_USER")
	QUEUE_HOST := viper.GetString("QUEUE_HOST")
	QUEUE_PASSWORD := viper.GetString("QUEUE_PASSWORD")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 🏓",
		})
	})
	queueStringConnection := fmt.Sprintf("amqp://%s:%s@%s", QUEUE_USER, QUEUE_PASSWORD, QUEUE_HOST)
	fmt.Println("cs >> "+queueStringConnection)
	conn, err := amqp.Dial(queueStringConnection)
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()
	r.Run(fmt.Sprintf(":%s", PORT)) 

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		USER_QUEUE, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

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
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}