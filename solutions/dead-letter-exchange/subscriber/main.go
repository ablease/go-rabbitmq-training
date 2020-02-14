package main

import (
	"log"

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

	err = ch.QueueBind(
		q.Name,                  // queue name
		"",                      // routing key
		"alive-letter-exchange", // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Reject messages to force onto dead letter queue
			err = d.Reject(false)
			util.FailOnError(err, "Failed to reject message")
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for dead-letter-messages. To exit press CTRL+C")
	<-forever
}
