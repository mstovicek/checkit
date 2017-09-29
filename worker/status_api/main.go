package main

import (
	"context"
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/messaging"
	"github.com/mstovicek/checkit/oauth"
	"github.com/mstovicek/checkit/repository_api"
	"github.com/mstovicek/checkit/repository_config"
)

type config struct {
	LoggerDebug              bool   `env:"LOGGER_DEBUG,default=true"`
	LoggerLogstashTcpAddress string `env:"LOGGER_LOGSTACH_TCP_ADDRESS,default=localhost:5000"`

	RepoConfigBasePath string `env:"REPO_CONFIG_BASE_PATH,default=./_configs_repo/"`
	RabbitMQURL        string `env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`
}

func main() {
	var cfg config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	log := logger.NewStdout().SetDebug(cfg.LoggerDebug).SetComponent("status-api")

	sendInspectionStatusConsumer, err := messaging.NewRabbitMqSendInspectionStatusConsumer(cfg.RabbitMQURL, log)
	if err != nil {
		panic(err)
	}

	defer sendInspectionStatusConsumer.Close()
	sentInspectionStatusChannel, err := sendInspectionStatusConsumer.Consume()
	if err != nil {
		panic(err)
	}

	sendPendingStatusConsumer, err := messaging.NewRabbitMqSentStatusPendingConsumer(cfg.RabbitMQURL, log)
	if err != nil {
		panic(err)
	}

	defer sendPendingStatusConsumer.Close()
	sendPendingStatusChannel, err := sendPendingStatusConsumer.Consume()
	if err != nil {
		panic(err)
	}

	repositoryConfigStore, err := repository_config.NewFileConfigStore(cfg.RepoConfigBasePath, log)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case delivery := <-sentInspectionStatusChannel:
			var message messaging.InspectionProcessed
			err := messaging.Unmarshal(delivery.Body, &message, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			if !repositoryConfigStore.HasConfig(message.Server, message.RepositoryName) {
				delivery.Nack(false, false)
				break
			}

			repoConfig, err := repositoryConfigStore.LoadConfig(message.Server, message.RepositoryName)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			auth, err := oauth.NewOAuth(repoConfig.Server, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			statusApi := repository_api.NewGithubStatusAPI(
				auth.Config.Client(context.Background(), &repoConfig.OAuthToken),
				message.RepositoryName,
				log,
			)
			err = statusApi.SetStatus(
				message.CommitHash,
				repository_api.CommitStatus{
					Status:      message.Status,
					DetailURL:   "http://89.221.208.88:8100/inspection/" + delivery.CorrelationId + "/",
					Description: "description up to 1k characters",
				},
			)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			log.Debug(logger.Fields{
				"correlationId": delivery.CorrelationId,
			}, "Sent inspection status")

			delivery.Ack(false)

		case delivery := <-sendPendingStatusChannel:
			var message messaging.CommitReceived
			err := messaging.Unmarshal(delivery.Body, &message, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			if !repositoryConfigStore.HasConfig(message.Server, message.RepositoryName) {
				delivery.Nack(false, false)
				break
			}

			repoConfig, err := repositoryConfigStore.LoadConfig(message.Server, message.RepositoryName)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			auth, err := oauth.NewOAuth(repoConfig.Server, log)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			statusApi := repository_api.NewGithubStatusAPI(
				auth.Config.Client(context.Background(), &repoConfig.OAuthToken),
				message.RepositoryName,
				log,
			)
			err = statusApi.SetStatus(
				message.CommitHash,
				repository_api.CommitStatus{
					Status:      repository_api.CommitStatusPending,
					Description: "description up to 1k characters",
				},
			)
			if err != nil {
				delivery.Nack(false, false)
				break
			}

			log.Debug(logger.Fields{
				"correlationId": delivery.CorrelationId,
			}, "Sent pending status")

			delivery.Ack(false)
		}
	}
}
