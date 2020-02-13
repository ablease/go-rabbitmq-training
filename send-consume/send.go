package main

import "time"

func main() {
	// TODO setup the connection with a ConnectionFactory (see slides for examples)

	// TODO setup a new channel using the newly created connection

	// TODO In the infinite loop while(true) add the following:
	// 1) wait for some period of time, call letsWait()
	// 2) get a quotation from the service by calling next()
	// 3) Send a message to the quotations exchange with routing key "nasq" by calling basicPublish() on the channel

	for {
	}
}

func letsWait() {
	time.Sleep(1000 * time.Millisecond)
}
