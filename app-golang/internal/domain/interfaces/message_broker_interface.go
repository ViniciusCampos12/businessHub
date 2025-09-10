package interfaces

type IMessageBroker interface {
	Publish(queue string, body []byte) error
}
