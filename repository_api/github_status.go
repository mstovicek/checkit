package repository_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mstovicek/checkit/logger"
	"net/http"
)

type statusJson struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

type githubStatusAPI struct {
	client         *http.Client
	repositoryName string
	statusMap      map[int]string
	logger         logger.Log
}

func NewGithubStatusAPI(client *http.Client, repositoryName string, logger logger.Log) StatusAPI {
	return &githubStatusAPI{
		client:         client,
		repositoryName: repositoryName,
		statusMap: map[int]string{
			CommitStatusPending: "pending",
			CommitStatusSuccess: "success",
			CommitStatusFailure: "failure",
			CommitStatusError:   "error",
		},
		logger: logger,
	}
}

func (api *githubStatusAPI) SetStatus(commitHash string, status CommitStatus) error {
	githubStatus := statusJson{
		State:       api.statusMap[status.Status],
		TargetURL:   status.DetailURL,
		Description: status.Description,
		Context:     apiGithubStatusContext,
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(githubStatus)

	url := ""
	url = fmt.Sprintf(
		// /repos/:owner/:repo/statuses/:sha
		"%s/repos/%s/statuses/%s",
		apiGithubEndpoint,
		api.repositoryName,
		commitHash,
	)

	api.logger.Info(logger.Fields{
		"url":  url,
		"body": body.String(),
	}, "Github set status sent")

	_, err := api.client.Post(
		url,
		"application/json",
		body,
	)
	if err != nil {
		api.logger.Error(logger.Fields{
			"url":   url,
			"body":  body.String(),
			"error": err.Error(),
		}, "Github set status could not be sent")
		return err
	}

	//defer resp.Body.Close()
	//
	//b, _ := httputil.DumpResponse(resp, true)
	//
	//fmt.Println(string(b))

	return nil
}
