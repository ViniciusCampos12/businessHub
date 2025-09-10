package interfaces

import amqp "github.com/rabbitmq/amqp091-go"

type IMessageBroker interface {
	Consume(queue string) (<-chan amqp.Delivery, error)
}
