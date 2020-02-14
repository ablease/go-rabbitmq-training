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

	err = ch.Tx()
	util.FailOnError(err, "Failed to Enable transactional channel")

	err = ch.ExchangeDeclare(
		"dead-letter-exchange", // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	dlq, err := ch.QueueDeclare(
		"dead-letter-queue", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	util.FailOnError(err, "Failed to declare dead letter queue")

	err = ch.QueueBind(
		dlq.Name,               // queue name
		"",                     // routing key
		"dead-letter-exchange", // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")

	err = ch.ExchangeDeclare(
		"alive-letter-exchange", // name
		"fanout",                // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	args := map[string]interface{}{
		"x-dead-letter-exchange": "dead-letter-exchange",
	}

	q, err := ch.QueueDeclare(
		"alive-letter-queue", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		args,                 // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	body := util.BodyFrom(os.Args)
	for {
		time.Sleep(1 * time.Second)

		log.Printf("Sending a message: %s", body)

		err = ch.Publish(
			"alive-letter-exchange", // exchange
			q.Name,                  // routing key
			false,                   // mandatory
			false,                   // immediate
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(body),
				DeliveryMode: amqp.Persistent,
			})
		util.FailOnError(err, "Failed to publish a message")

		err = ch.TxCommit()
		util.FailOnError(err, "Failed to commit a commit")

	}
}
