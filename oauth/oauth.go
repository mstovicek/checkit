package oauth

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/mstovicek/checkit/logger"
	"golang.org/x/oauth2"
	"math/rand"
)

type OAuth struct {
	log    logger.Log
	Config *oauth2.Config
}

func NewOAuth(server string, log logger.Log) (*OAuth, error) {
	config := newOAuthConfig(server)
	if config == nil {
		log.Error(logger.Fields{
			"server": server,
		}, "cannot load OAuth config")
		return nil, errors.New("cannot load OAuth config")
	}

	return &OAuth{
		log:    log,
		Config: config,
	}, nil
}

func newOAuthConfig(server string) *oauth2.Config {
	var config *oauth2.Config

	switch server {
	case GithubServer:
		config = &githubConfig
	default:
		config = nil
	}

	return config
}

func (auth *OAuth) ExchangeAccessToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()

	token, err := auth.Config.Exchange(ctx, code)
	if err != nil {
		auth.log.Error(logger.Fields{}, "cannot exchange token")
		return nil, err
	}

	return token, nil
}

func (auth *OAuth) GetAuthUrl(state string) string {
	return auth.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (auth *OAuth) GetRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
