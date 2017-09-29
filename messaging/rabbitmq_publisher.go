package messaging

import (
	"encoding/json"
	"github.com/mstovicek/checkit/logger"
	"github.com/streadway/amqp"
)

const (
	RabbitMqURL = "amqp://guest:guest@localhost:5672/"

	RabbitMqExchangeName = "checkit"

	RabbitMqTopicCommitReceived      = "commit.received"
	RabbitMqTopicInspectionProcessed = "inspection.processed"
)

type rabbitMqPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	topic      string
	log        logger.Log
}

func NewRabbitMqCommitReceivedPublisher(queueURL string, log logger.Log) (QueuePublisher, error) {
	return NewRabbitMqPublisher(
		queueURL,
		RabbitMqTopicCommitReceived,
		log,
	)
}

func NewRabbitMqInspectionProcessedPublisher(queueURL string, log logger.Log) (QueuePublisher, error) {
	return NewRabbitMqPublisher(
		queueURL,
		RabbitMqTopicInspectionProcessed,
		log,
	)
}

func NewRabbitMqPublisher(queueURL string, topic string, log logger.Log) (QueuePublisher, error) {
	conn, err := amqp.Dial(queueURL)
	if err != nil {
		log.Error(logger.Fields{
			"queueURL": queueURL,
			"error":    err.Error(),
		}, "Failed to connect to RabbitMQ")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error(logger.Fields{
			"error": err.Error(),
		}, "Failed to open a channel")
		return nil, err
	}

	err = ch.ExchangeDeclare(
		RabbitMqExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error(logger.Fields{
			"error": err.Error(),
		}, "Failed to declare an exchange")
		return nil, err
	}

	return &rabbitMqPublisher{
		connection: conn,
		channel:    ch,
		topic:      topic,
		log:        log,
	}, nil
}

func (rabbitMqPublisher *rabbitMqPublisher) Close() {
	rabbitMqPublisher.connection.Close()
	rabbitMqPublisher.channel.Close()
}

func (rabbitMqPublisher *rabbitMqPublisher) Publish(correlationId string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		rabbitMqPublisher.log.Error(logger.Fields{
			"correlationId": correlationId,
			"message":       message,
			"error":         err.Error(),
		}, "Cannot marshal a message")
		return err
	}

	rabbitMqPublisher.log.Debug(logger.Fields{
		"topic":         rabbitMqPublisher.topic,
		"correlationId": correlationId,
		"body":          string(body),
	}, "publish message")

	err = rabbitMqPublisher.channel.Publish(
		RabbitMqExchangeName,
		rabbitMqPublisher.topic,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          body,
		},
	)
	if err != nil {
		rabbitMqPublisher.log.Error(logger.Fields{
			"correlationId": correlationId,
			"message":       message,
			"error":         err.Error(),
		}, "Failed to publish a message")
		return err
	}

	return nil
}
