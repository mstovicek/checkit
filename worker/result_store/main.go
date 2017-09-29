package main

import (
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/checkit/document_store"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/messaging"
)

type config struct {
	LoggerDebug              bool   `env:"LOGGER_DEBUG,default=true"`
	LoggerLogstashTcpAddress string `env:"LOGGER_LOGSTACH_TCP_ADDRESS,default=localhost:5000"`

	RabbitMQURL string `env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`

	MongoDBURL            string `env:"MONGODB_URL,default=mongodb://localhost:27017"`
	MongoDBDatabaseName   string `env:"MONGODB_DB_NAME,default=checkit"`
	MongoDBCollectionName string `env:"MONGODB_COLLECTION_NAME,default=inspection_results"`
}

func main() {
	var cfg config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	log := logger.NewStdout().SetDebug(cfg.LoggerDebug).SetComponent("result-store")

	storeInspectionResultConsumer, err := messaging.NewRabbitMqStoreInspectionResult(cfg.RabbitMQURL, log)
	if err != nil {
		panic(err)
	}

	defer storeInspectionResultConsumer.Close()
	storeInspectionResultChannel, err := storeInspectionResultConsumer.Consume()
	if err != nil {
		panic(err)
	}

	documentStore, err := document_store.NewMongoDB(log, cfg.MongoDBURL, cfg.MongoDBDatabaseName, cfg.MongoDBCollectionName)
	if err != nil {
		panic("cannot connect to MongoDB")
	}

	for {
		select {
		case delivery := <-storeInspectionResultChannel:
			var message messaging.InspectionProcessed
			err := messaging.Unmarshal(delivery.Body, &message, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			document := document_store.InspectionResultDocument{
				Uuid:           delivery.CorrelationId,
				Time:           message.Time,
				Server:         message.Server,
				RepositoryName: message.RepositoryName,
				CommitHash:     message.CommitHash,
				Status:         message.Status,
				FixedFiles:     message.FixedFiles,
			}

			err = documentStore.Insert(&document)
			if err != nil {
				delivery.Nack(false, true)
				break
			}

			log.Debug(logger.Fields{
				"correlationId": delivery.CorrelationId,
			}, "Stored document")

			delivery.Ack(false)
		}
	}
}
