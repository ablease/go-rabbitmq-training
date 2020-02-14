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

	noWait := false
	err = ch.Confirm(noWait)
	util.FailOnError(err, "Failed to enable publisher confirms")

	err = ch.ExchangeDeclare(
		"quorum-exchange", // name
		"fanout",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	args := map[string]interface{}{
		"x-queue-type": "quorum",
	}

	q, err := ch.QueueDeclare(
		"quorum-queue", // name
		true,           // durable
		false,          // autodelete
		false,          // exclusive
		false,          // nowait
		args,           // args
	)
	util.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,            // queue name
		"",                // routing key
		"quorum-exchange", // exchange
		false,             // nowait
		nil,               // args
	)

	util.FailOnError(err, "Failed to bind a queue")

	body := util.BodyFrom(os.Args)
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	for {
		time.Sleep(1 * time.Second)
		err = ch.Publish(
			"quorum-exchange", // exchange
			"",                // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)

		confirmOne(confirms)

	}
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	// receive a confirm from the confirms channel, if the confirmed.Ack is true then...
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
