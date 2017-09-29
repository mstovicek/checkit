package api

import (
	"encoding/json"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/messaging"
	"github.com/mstovicek/checkit/oauth"
	"io"
	"net/http"
)

type githubPushJson struct {
	HeadCommit struct {
		CommitHash    string   `json:"id"`
		ModifiedFiles []string `json:"modified"`
	} `json:"head_commit"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type webhookGithubHandler struct {
	logger        logger.Log
	publisher     messaging.QueuePublisher
	uuidGenerator messaging.UUIDGenerator
}

func newWebhookGithubHandler(log logger.Log, publisher messaging.QueuePublisher, uuidGen messaging.UUIDGenerator) http.Handler {
	return &webhookGithubHandler{
		logger:        log,
		publisher:     publisher,
		uuidGenerator: uuidGen,
	}
}

func (h *webhookGithubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pushJson githubPushJson
	err := decoder.Decode(&pushJson)
	if err != nil {
		h.logger.Errorf("Error decoding push hook, %s\n", err)
	}

	defer r.Body.Close()

	err = h.publisher.Publish(
		h.uuidGenerator.Generate(),
		messaging.CommitReceived{
			CommitHash:     pushJson.HeadCommit.CommitHash,
			Server:         oauth.GithubServer,
			RepositoryName: pushJson.Repository.FullName,
			Files:          pushJson.HeadCommit.ModifiedFiles,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, "{}")
}
