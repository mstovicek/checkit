package messaging

import (
	"github.com/mstovicek/checkit/logger"
	"github.com/streadway/amqp"
)

const (
	RabbitMqQueueRunInspection         = "runInspection"
	RabbitMqQueueSendStatusPending     = "sendStatusPending"
	RabbitMqQueueSendInspectionStatus  = "sendInspectionStatus"
	RabbitMqQueueStoreInspectionResult = "storeInspectionResult"
)

type rabbitMqConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	log        logger.Log
}

func NewRabbitMqRunInspectionConsumer(queueURL string, log logger.Log) (QueueConsumer, error) {
	return NewRabbitMqConsumer(
		queueURL,
		RabbitMqExchangeName,
		RabbitMqQueueRunInspection,
		RabbitMqTopicCommitReceived,
		1,
		log,
	)
}

func NewRabbitMqSentStatusPendingConsumer(queueURL string, log logger.Log) (QueueConsumer, error) {
	return NewRabbitMqConsumer(
		queueURL,
		RabbitMqExchangeName,
		RabbitMqQueueSendStatusPending,
		RabbitMqTopicCommitReceived,
		1,
		log,
	)
}

func NewRabbitMqSendInspectionStatusConsumer(queueURL string, log logger.Log) (QueueConsumer, error) {
	return NewRabbitMqConsumer(
		queueURL,
		RabbitMqExchangeName,
		RabbitMqQueueSendInspectionStatus,
		RabbitMqTopicInspectionProcessed,
		1,
		log,
	)
}

func NewRabbitMqStoreInspectionResult(queueURL string, log logger.Log) (QueueConsumer, error) {
	return NewRabbitMqConsumer(
		queueURL,
		RabbitMqExchangeName,
		RabbitMqQueueStoreInspectionResult,
		RabbitMqTopicInspectionProcessed,
		1,
		log,
	)
}

func NewRabbitMqConsumer(queueUrl string, exchangeName string, queueName string, topic string, prefetchCount int, log logger.Log) (QueueConsumer, error) {
	conn, err := amqp.Dial(queueUrl)
	if err != nil {
		log.Error(logger.Fields{
			"queueUrl": queueUrl,
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
		exchangeName,
		amqp.ExchangeTopic,
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

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error(logger.Fields{
			"error": err.Error(),
		}, "Failed to declare a queue")
		return nil, err
	}

	err = ch.QueueBind(
		queueName,
		topic,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Error(logger.Fields{
			"error": err.Error(),
		}, "Failed to bind a queue")
		return nil, err
	}

	err = ch.Qos(
		prefetchCount,
		0,
		false,
	)
	if err != nil {
		log.Error(logger.Fields{
			"error": err.Error(),
		}, "Failed to set QoS")
		return nil, err
	}

	return &rabbitMqConsumer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
		log:        log,
	}, nil
}

func (rabbitMqConsumer *rabbitMqConsumer) Close() {
	rabbitMqConsumer.channel.Close()
	rabbitMqConsumer.connection.Close()
}

func (rabbitMqConsumer *rabbitMqConsumer) Consume() (<-chan amqp.Delivery, error) {
	messages, err := rabbitMqConsumer.channel.Consume(
		rabbitMqConsumer.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		rabbitMqConsumer.log.Error(logger.Fields{
			"error": err.Error(),
		}, "Failed to register a consumer")
		return nil, err
	}

	rabbitMqConsumer.log.Info(logger.Fields{
		"queue": rabbitMqConsumer.queueName,
	}, "Consuming messages")

	return messages, nil
}
