package main

import (
	"github.com/streadway/amqp"
)

func main() {
	// Create a connection to your local RabbitMQ.
	_, _ = amqp.Dial("")

	// Create a Channel from the Connection.

	// Declare a Queue.

	// Call Publish on the Channel to send 1 message to the queue
}
