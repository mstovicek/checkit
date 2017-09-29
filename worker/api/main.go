package main

import (
	"github.com/joeshaw/envdecode"
	"github.com/mstovicek/checkit/api"
	"github.com/mstovicek/checkit/logger"
)

type config struct {
	Listen string `env:"LISTEN,default=:8100"`

	LoggerDebug              bool   `env:"LOGGER_DEBUG,default=true"`
	LoggerLogstashTcpAddress string `env:"LOGGER_LOGSTACH_TCP_ADDRESS,default=localhost:5000"`

	BaseAuthURL       string `env:"BASE_URL_AUTH,default=/auth/"`
	BaseWebhookURL    string `env:"BASE_URL_WEBHOOK,default=/webhook/"`
	BaseRepositoryURL string `env:"BASE_URL_REPOSITORY,default=/repository/"`
	BaseInspectionURL string `env:"BASE_URL_INSPECTION,default=/inspection/"`

	SessionAuthenticationKey string `env:"SESSION_AUTHENTICATION_KEY,required"`
	SessionEncryptionKey     string `env:"SESSION_ENCRYPTION_KEY"`
	SessionStoreName         string `env:"SESSION_STORE_NAME,default=OAuthGithub"`

	RabbitMQURL string `env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`

	MongoDBURL            string `env:"MONGODB_URL,default=mongodb://localhost:27017"`
	MongoDBDatabaseName   string `env:"MONGODB_DB_NAME,default=checkit"`
	MongoDBCollectionName string `env:"MONGODB_COLLECTION_NAME,default=inspection_results"`

	RepoConfigBasePath string `env:"REPO_CONFIG_BASE_PATH,default=./_configs_repo/"`
}

func main() {
	var cfg config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	log := logger.NewStdout().SetDebug(cfg.LoggerDebug).SetComponent("api")

	server := api.NewServer(cfg.Listen, log)

	apiHandlers, err := api.GetOAuthHandlers(
		cfg.BaseAuthURL,
		log,
		cfg.SessionAuthenticationKey,
		cfg.SessionEncryptionKey,
		cfg.SessionStoreName,
	)
	if err != nil {
		panic(err)
	}
	server.AddHandlers(apiHandlers)

	webookHandlers, err := api.GetWebhookHandlers(
		cfg.BaseWebhookURL,
		log,
		cfg.RabbitMQURL,
	)
	if err != nil {
		panic(err)
	}
	server.AddHandlers(webookHandlers)

	repositoryHandlers, err := api.GetRepositoryHandlers(
		cfg.BaseRepositoryURL,
		log,
		cfg.SessionAuthenticationKey,
		cfg.SessionEncryptionKey,
		cfg.SessionStoreName,
		cfg.RepoConfigBasePath,
	)
	if err != nil {
		panic(err)
	}
	server.AddHandlers(repositoryHandlers)

	inspectionHandlers, err := api.GetInspectionHandlers(
		cfg.BaseInspectionURL,
		log,
		cfg.SessionAuthenticationKey,
		cfg.SessionEncryptionKey,
		cfg.SessionStoreName,
		cfg.MongoDBURL,
		cfg.MongoDBDatabaseName,
		cfg.MongoDBCollectionName,
	)
	if err != nil {
		panic(err)
	}
	server.AddHandlers(inspectionHandlers)

	server.Run()
}
