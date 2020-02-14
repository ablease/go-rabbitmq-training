package main

import (
	"github.com/streadway/amqp"
)

func main() {
	// Dial RabbitMQ and create a Channel
	_, _ = amqp.Dial("")

	// Declare an Exchange for our publisher/subscriber

	// Declare a queue for our consumer

	// Bind the Queue we declared to our Exchange

	// Loop forever and consume messages from the queue
}
