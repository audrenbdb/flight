package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"laresa"
	"laresa/ulid"
	"log"
	"net/http"
	"os"
)

func main() {
	amqpURL := "amqp://guest:guest@localhost:5672/"
	serviceAddr := "127.0.0.1:8080"

	service := &laresa.Booker{
		NewID:     ulid.NewBuilder(),
		Repo:      laresa.NewInMemRepo(),
		Publisher: laresa.NewAMQPPublisher(amqpURL),
	}
	handler := laresa.NewChiHTTPHandler(service)
	go http.ListenAndServe(serviceAddr, handler)
	fmt.Printf("Server listening on: %s\n", serviceAddr)
	newAMQPListener(amqpURL)
}

// newAMQPListener start a small listener to register messages into a file
// it is the "micro service" from Valentina's challenge :p
func newAMQPListener(url string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"reservation", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	for d := range msgs {
		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile("reservations", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(append(d.Body, []byte("\n")...)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
