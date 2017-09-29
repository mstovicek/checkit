package api

import (
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"net/http"
)

func GetRepositoryHandlers(
	baseUrl string,
	logger logger.Log,
	sessionAuthenticationKey string,
	sessionEncryptionKey string,
	sessionStoreName string,
	repoConfigBasePath string,
) ([]Handler, error) {
	sessionStore := sessions.NewCookieStore(
		[]byte(sessionAuthenticationKey),
		[]byte(sessionEncryptionKey),
	)

	return []Handler{
		{
			path: baseUrl,
			handler: newHasSessionMiddleware(
				logger,
				sessionStore,
				sessionStoreName,
				newRepositoryIndexHandler(
					logger,
					sessionStore,
					sessionStoreName,
					repoConfigBasePath,
				),
			),
			methods: []string{http.MethodGet},
		},
		{
			path: baseUrl + "add/",
			handler: newHasSessionMiddleware(
				logger,

				sessionStore,
				sessionStoreName,
				newAddRepositoryHandler(
					logger,
					sessionStore,
					sessionStoreName,
					repoConfigBasePath,
				),
			),
			methods: []string{http.MethodGet},
		},
	}, nil
}
