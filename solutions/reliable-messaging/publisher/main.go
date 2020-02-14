package main

import (
	"log"
	"os"
	"time"

	"github.com/ablease/rabbitmq-training/util"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"reliable-message-exchange", // name
		"fanout",                    // type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		nil,                         // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"reliable-message-queue", // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	body := util.BodyFrom(os.Args)
	for {
		time.Sleep(1 * time.Second)

		log.Printf("Sending a message: %s", body)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(body),
				DeliveryMode: amqp.Persistent,
			})
		util.FailOnError(err, "Failed to publish a message")

	}
}
