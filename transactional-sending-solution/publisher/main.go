package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Tx()
	failOnError(err, "Failed to Enable transactional channel")

	err = ch.ExchangeDeclare(
		"transactional-sending-exchange", // name
		"fanout",                         // type
		true,                             // durable
		false,                            // auto-deleted
		false,                            // internal
		false,                            // no-wait
		nil,                              // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"transactional-sending-queue", // name
		true,                          // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := bodyFrom(os.Args)
	for {
		time.Sleep(1 * time.Second)

		log.Printf("Sending a message: %s", body)

		err = ch.Publish(
			"transactional-sending-exchange", // exchange
			q.Name,                           // routing key
			false,                            // mandatory
			false,                            // immediate
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(body),
				DeliveryMode: amqp.Persistent,
			})
		failOnError(err, "Failed to publish a message")

		err = ch.TxCommit()
		failOnError(err, "Failed to commit a commit")

	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
