package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connection *amqp.Connection
}

func (p *Producer) setup() error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (p *Producer) Push(event string, severity string) error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func NewEventProducer(conn *amqp.Connection) (Producer, error) {
	producer := Producer{
		connection: conn,
	}

	err := producer.setup()
	if err != nil {
		return Producer{}, err
	}

	return producer, nil
}
