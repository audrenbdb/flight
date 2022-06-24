package laresa

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type inMemPublisher struct{}

func NewInMemPublisher() *inMemPublisher {
	return &inMemPublisher{}
}

func (p *inMemPublisher) PublishReservation(ctx context.Context, resa Reservation) error {
	log.Printf("New reservation published: %#v\n", resa)
	return nil
}

const amqpQueueName = "reservation"

type amqpSender struct {
	*amqp.Connection
	*amqp.Channel
}

func NewAMQPPublisher(url string) *amqpSender {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	_, err = ch.QueueDeclare(
		amqpQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return &amqpSender{
		Connection: conn,
		Channel:    ch,
	}
}

func (s *amqpSender) PublishReservation(ctx context.Context, resa Reservation) error {
	b, err := json.Marshal(resa)
	if err != nil {
		return err
	}
	return s.Publish("", amqpQueueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        b,
	})
}
