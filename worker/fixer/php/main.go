package main

import (
	"encoding/json"
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/checkit/code_fixer"
	"github.com/mstovicek/checkit/git"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/messaging"
	"github.com/mstovicek/checkit/repository_config"
	"time"
)

type config struct {
	LoggerDebug              bool   `env:"LOGGER_DEBUG,default=true"`
	LoggerLogstashTcpAddress string `env:"LOGGER_LOGSTACH_TCP_ADDRESS,default=localhost:5000"`

	RepositoryConfigBasePath  string `env:"REPO_CONFIG_BASE_PATH,default=./_configs_repo/"`
	RepositoryStorageBasePath string `env:"REPO_BASE_PATH,default=./_repo/"`
	RabbitMQURL               string `env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`

	CommandName string `env:"COMMAND_NAME,default=./php-cs-fixer"`
}

func main() {
	var cfg config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	log := logger.NewStdout().SetDebug(cfg.LoggerDebug).SetComponent("fixer-php")

	consumer, err := messaging.NewRabbitMqRunInspectionConsumer(cfg.RabbitMQURL, log)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	runInspectionChannel, err := consumer.Consume()
	if err != nil {
		panic(err)
	}

	producer, err := messaging.NewRabbitMqInspectionProcessedPublisher(cfg.RabbitMQURL, log)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	for {
		select {
		case delivery := <-runInspectionChannel:
			var message messaging.CommitReceived
			err := messaging.Unmarshal(delivery.Body, &message, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			configStore, err := repository_config.NewFileConfigStore(cfg.RepositoryConfigBasePath, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			if !configStore.HasConfig(message.Server, message.RepositoryName) {
				delivery.Nack(false, false)
				break
			}

			config, err := configStore.LoadConfig(message.Server, message.RepositoryName)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			client := git.NewFilesystemClient(cfg.RepositoryStorageBasePath, log)
			directory, err := client.Checkout(config.CloneUrl, config.Server, config.RepositoryName, message.CommitHash)
			if err != nil {
				log.Error(logger.Fields{
					"cloneUrl":       config.CloneUrl,
					"server":         config.Server,
					"repositoryName": config.RepositoryName,
					"error":          err.Error(),
				}, "cannot checkout")
				delivery.Nack(false, false)
				break
			}

			fixer := code_fixer.NewPhp(log, code_fixer.PhpConfiguration{}, cfg.CommandName)
			fixerResult, err := fixer.Run(directory, message.Files)
			if err != nil {
				log.Error(logger.Fields{}, "cannot unmarshal fixer result")
				delivery.Nack(false, false)
				break
			}

			b, err := json.MarshalIndent(fixerResult, "", "  ")
			if err != nil {
				log.Error(logger.Fields{}, "cannot unmarshal fixer result")
				delivery.Nack(false, false)
				break
			}

			log.Debug(logger.Fields{
				"message":       message,
				"email":         config.Email,
				"correlationId": delivery.CorrelationId,
				"repo dir":      directory,
				"result":        string(b),
			}, "Inspection done")

			err = producer.Publish(
				delivery.CorrelationId,
				messaging.InspectionProcessed{
					Time:           time.Now(),
					Server:         message.Server,
					RepositoryName: message.RepositoryName,
					CommitHash:     message.CommitHash,
					Status:         fixerResult.Status,
					FixedFiles:     fixerResult.FixedFiles,
				},
			)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			delivery.Ack(false)
		}
	}
}
