package api

import (
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/messaging"
	"github.com/mstovicek/checkit/oauth"
	"net/http"
)

func GetWebhookHandlers(
	baseUrl string,
	logger logger.Log,
	rabbitMQURL string,
) ([]Handler, error) {
	commitReceivedPublisher, err := messaging.NewRabbitMqCommitReceivedPublisher(rabbitMQURL, logger)
	if err != nil {
		return []Handler{}, err
	}

	return []Handler{
		{
			path: baseUrl + oauth.GithubServer + "/",
			handler: newWebhookGithubHandler(
				logger,
				commitReceivedPublisher,
				messaging.NewUUIDGenerator(),
			),
			methods: []string{http.MethodPost},
		},
	}, nil
}
