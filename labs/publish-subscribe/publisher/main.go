package main

import (
	"github.com/streadway/amqp"
)

func main() {
	// Dial RabbitMQ and create a connection
	_, _ = amqp.Dial("")

	// Declare an exchange for our publisher/subscriber

	// Loop and wait for a small time, then publish a message to the Exchange
}
