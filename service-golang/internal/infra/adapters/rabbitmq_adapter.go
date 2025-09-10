package adapters

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMqAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMqAdapter(url string) *RabbitMqAdapter {
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	return &RabbitMqAdapter{
		conn:    conn,
		channel: ch,
	}
}

func (r *RabbitMqAdapter) Consume(queue string) (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to rabbitmq queue declare: %w", err)
	}

	data, err := r.channel.Consume(
		queue,
		"",    // consumer
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)

	if err != nil {
		log.Fatalf("fail to load consumer: %v", err)
	}

	return data, nil
}

func (r *RabbitMqAdapter) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("error closing channel: %w", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}
	return nil
}
