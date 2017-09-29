package repository_api

import (
	"errors"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"net/http"
)

const (
	CommitStatusPending = 0
	CommitStatusSuccess = 1
	CommitStatusFailure = 2
	CommitStatusError   = 3
)

type UserAPI interface {
	GetUser() (User, error)
}

type StatusAPI interface {
	SetStatus(commitHash string, status CommitStatus) error
}

type RepositoryAPI interface {
	GetRepositories() ([]Repository, error)
	GetRepository(fullName string) (Repository, error)
}

func NewRepositoryAPI(server string, client *http.Client, log logger.Log) (RepositoryAPI, error) {
	switch server {
	case oauth.GithubServer:
		return newGithubRepositoryAPI(client, log), nil
	default:
		log.Error(logger.Fields{
			"server": server,
		}, "invalid repository server")
		return nil, errors.New("invalid repository server")
	}
}

type CommitStatus struct {
	Status      int
	DetailURL   string
	Description string
}

type User struct {
	Name  string
	Email string
}

type Repository struct {
	Server   string
	FullName string
	URL      string
	CloneURL string
	IsAdmin  bool
	CanRead  bool
	CanWrite bool
}
