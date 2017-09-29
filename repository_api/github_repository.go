package repository_api

import (
	"encoding/json"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"io/ioutil"
	"net/http"
)

type githubRepositoryJson struct {
	FullName    string `json:"full_name"`
	HtmlURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	Permissions struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

type githubRepositoryAPI struct {
	client *http.Client
	logger logger.Log
}

func newGithubRepositoryAPI(client *http.Client, logger logger.Log) RepositoryAPI {
	return &githubRepositoryAPI{
		client: client,
		logger: logger,
	}
}

func (api *githubRepositoryAPI) GetRepositories() ([]Repository, error) {
	endpoint := apiGithubEndpoint + "/user/repos"
	resp, err := api.client.Get(endpoint)
	if err != nil {
		api.logger.Error(logger.Fields{
			"endpoint": endpoint,
			"error":    err.Error(),
		}, "Cannot fetch repositories")
		return []Repository{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.logger.Error(logger.Fields{
			"endpoint": endpoint,
			"body":     body,
			"error":    err.Error(),
		}, "Cannot read repositories body")
		return []Repository{}, err
	}

	var repositories []githubRepositoryJson
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		api.logger.Error(logger.Fields{
			"endpoint": endpoint,
			"body":     body,
			"error":    err.Error(),
		}, "Cannot unmarshal repositories response")
		return []Repository{}, err
	}

	var output []Repository
	for _, repository := range repositories {
		output = append(
			output,
			Repository{
				Server:   oauth.GithubServer,
				FullName: repository.FullName,
				URL:      repository.HtmlURL,
				CloneURL: repository.CloneURL,
				IsAdmin:  repository.Permissions.Admin,
				CanRead:  repository.Permissions.Pull,
				CanWrite: repository.Permissions.Push,
			},
		)
	}

	return output, nil
}

func (api *githubRepositoryAPI) GetRepository(fullName string) (Repository, error) {
	resp, err := api.client.Get(apiGithubEndpoint + "/repos/" + fullName)
	if err != nil {
		api.logger.Error(logger.Fields{
			"fullName": fullName,
			"error":    err.Error(),
		}, "Cannot fetch repository")
		return Repository{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.logger.Error(logger.Fields{
			"fullName": fullName,
			"body":     body,
			"error":    err.Error(),
		}, "Cannot read repository body")
		return Repository{}, err
	}

	var repository githubRepositoryJson
	err = json.Unmarshal(body, &repository)
	if err != nil {
		api.logger.Error(logger.Fields{
			"fullName": fullName,
			"body":     body,
			"error":    err.Error(),
		}, "Cannot unmarshal repository response")
		return Repository{}, err
	}

	return Repository{
		Server:   oauth.GithubServer,
		FullName: repository.FullName,
		URL:      repository.HtmlURL,
		CloneURL: repository.CloneURL,
		IsAdmin:  repository.Permissions.Admin,
		CanRead:  repository.Permissions.Pull,
		CanWrite: repository.Permissions.Push,
	}, nil
}
