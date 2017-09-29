package api

import (
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"net/http"
)

const (
	sessionKeyAccessToken = "accessToken"
	sessionKeyEmail       = "email"
	sessionKeyOAuthServer = "server"
	sessionKeyState       = "state"
	queryKeyState         = "state"
)

func GetOAuthHandlers(
	baseUrl string,
	logger logger.Log,
	sessionAuthenticationKey string,
	sessionEncryptionKey string,
	sessionStoreName string,
) ([]Handler, error) {
	sessionStore := sessions.NewCookieStore(
		[]byte(sessionAuthenticationKey),
		[]byte(sessionEncryptionKey),
	)

	return []Handler{
		{
			path: baseUrl,
			handler: newAuthIndexHandler(
				logger,
				sessionStore,
				sessionStoreName,
			),
			methods: []string{http.MethodGet},
		},
		{
			path: baseUrl + "logout/",
			handler: newAuthLogoutHandler(
				logger,
				sessionStore,
				sessionStoreName,
			),
			methods: []string{http.MethodGet},
		},
		{
			path: baseUrl + oauth.GithubServer + "/",
			handler: newAuthGithubHandler(
				logger,
				sessionStore,
				sessionStoreName,
			),
			methods: []string{http.MethodGet},
		},
	}, nil
}
