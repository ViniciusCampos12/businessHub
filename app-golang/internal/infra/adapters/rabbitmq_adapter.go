package adapters

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
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

func (r *RabbitMqAdapter) Publish(queue string, body []byte) error {
	_, err := r.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to rabbitmq queue declare: %w", err)
	}

	return r.channel.Publish(
		"",    // exchange
		queue, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
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
