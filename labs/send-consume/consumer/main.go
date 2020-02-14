package main

import (
	"github.com/streadway/amqp"
)

func main() {

	// Create a connection to your local RabbitMQ.
	_, _ = amqp.Dial("")

	// Create a Channel from the Connection.

	// Declare a Queue. Tip: If you declare a queue that doesnt exist,
	// you might never receive any messages

	// Consume a message from the queue.

}
