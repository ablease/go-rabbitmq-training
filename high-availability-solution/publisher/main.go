package main

import (
	"log"
	"os"
	"time"

	"github.com/ablease/rabbitmq-training/util"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"high-availability-exchange", // name
		"fanout",                     // type
		true,                         // durable
		false,                        // auto-deleted
		false,                        // internal
		false,                        // no-wait
		nil,                          // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"high-availability-queue", // name
		true,                      // durable
		false,                     // autodelete
		false,                     // exclusive
		false,                     // nowait
		nil,                       // args
	)
	util.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                       // queue name
		"",                           // routing key
		"high-availability-exchange", // exchange
		false,                        // nowait
		nil,                          // args
	)

	util.FailOnError(err, "Failed to bind a queue")

	body := util.BodyFrom(os.Args)

	for {
		time.Sleep(1 * time.Second)
		err = ch.Publish(
			"high-availability-exchange", // exchange
			"",                           // routing key
			false,                        // mandatory
			false,                        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}
}
