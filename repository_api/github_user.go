package repository_api

import (
	"encoding/json"
	"errors"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"io/ioutil"
	"net/http"
)

type userJson struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type emailJson struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type githubUserAPI struct {
	client *http.Client
	logger logger.Log
}

func NewGithubUserAPI(client *http.Client, log logger.Log) UserAPI {
	return &githubUserAPI{
		client: client,
		logger: log,
	}
}

func (api *githubUserAPI) GetUser() (User, error) {
	endpoint := apiGithubEndpoint + "/user"
	resp, err := api.client.Get(endpoint)
	if err != nil {
		api.logger.Error(logger.Fields{
			"server":   oauth.GithubServer,
			"endpoint": endpoint,
		}, "cannot get user")
		return User{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.logger.Error(logger.Fields{
			"server":   oauth.GithubServer,
			"endpoint": endpoint,
		}, "cannot read user response")
		return User{}, err
	}

	var user userJson
	err = json.Unmarshal(body, &user)
	if err != nil {
		api.logger.Error(logger.Fields{
			"server":   oauth.GithubServer,
			"endpoint": endpoint,
		}, "cannot unmarshal user response")
		return User{}, err
	}

	if user.Email == "" {
		user.Email, _ = api.getEmail()
	}

	if user.Email == "" {
		api.logger.Error(logger.Fields{
			"endpoint": endpoint,
		}, "cannot find any email")
		return User{}, errors.New("cannot find any email")
	}

	return User{
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (api *githubUserAPI) getEmail() (string, error) {
	endpoint := apiGithubEndpoint + "/user/emails"
	resp, err := api.client.Get(endpoint)
	if err != nil {
		api.logger.Error(logger.Fields{
			"server":   oauth.GithubServer,
			"endpoint": endpoint,
		}, "cannot get email")
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.logger.Error(logger.Fields{
			"server":   oauth.GithubServer,
			"endpoint": endpoint,
		}, "cannot read email response")
		return "", err
	}

	var emails []emailJson
	err = json.Unmarshal(body, &emails)
	if err != nil {
		api.logger.Error(logger.Fields{
			"server":   oauth.GithubServer,
			"endpoint": endpoint,
		}, "cannot unmarshal email response")
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	api.logger.Error(logger.Fields{
		"server":   oauth.GithubServer,
		"endpoint": endpoint,
	}, "cannot find email")

	return "", errors.New("cannot find email")
}
