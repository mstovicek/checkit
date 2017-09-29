package messaging

import (
	"github.com/streadway/amqp"
)

type QueueConsumer interface {
	Close()
	Consume() (<-chan amqp.Delivery, error)
}

type QueuePublisher interface {
	Close()
	Publish(correlationId string, message interface{}) error
}
