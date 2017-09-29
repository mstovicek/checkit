package api

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"github.com/mstovicek/checkit/repository_api"
	"github.com/mstovicek/checkit/repository_config"
	"io"
	"net/http"
	"net/url"
)

type repositoryIndexHandler struct {
	logger             logger.Log
	sessionStore       sessions.Store
	sessionStoreName   string
	repoConfigBasePath string
}

func newRepositoryIndexHandler(
	logger logger.Log,
	sessionStore sessions.Store,
	sessionStoreName string,
	repoConfigBasePath string,
) http.Handler {
	return &repositoryIndexHandler{
		logger:             logger,
		sessionStore:       sessionStore,
		sessionStoreName:   sessionStoreName,
		repoConfigBasePath: repoConfigBasePath,
	}
}

func (h *repositoryIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionStoreName)
	if err != nil {
		h.logger.Error(logger.Fields{
			"sessionStoreName": h.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot get session store")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	server, ok := session.Values[sessionKeyOAuthServer].(string)
	if !ok {
		if err != nil {
			writeError(w, http.StatusBadRequest, "Unsupported OAuth server")
			return
		}
	}

	auth, err := oauth.NewOAuth(server, h.logger)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Unsupported OAuth server")
		return
	}

	accessToken, err := oauth.DecodeToken(session.Values[sessionKeyAccessToken])
	if err != nil {
		h.logger.Error(logger.Fields{
			"server": server,
			"error":  err.Error(),
		}, "Cannot decode access token")
		writeError(w, http.StatusInternalServerError, "Cannot decode access token")
		return
	}

	repositoryApi, err := repository_api.NewRepositoryAPI(
		server,
		auth.Config.Client(context.Background(), accessToken),
		h.logger,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Cannot connect to repository API")
		return
	}

	repos, err := repositoryApi.GetRepositories()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot fetch repositories")
		return
	}

	repositoriesHtml := ""
	for _, repo := range repos {
		if !repo.IsAdmin {
			continue
		}

		configStore, err := repository_config.NewFileConfigStore(h.repoConfigBasePath, h.logger)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "cannot initialize repository config store")
			continue
		}

		var action string
		if configStore.HasConfig(server, repo.FullName) {
			action = "added"
		} else {
			action = "<a href=/repository/add/?repository=" + url.QueryEscape(repo.FullName) + ">add</a>"
		}

		repositoriesHtml += fmt.Sprintf(
			"<li>%s: %s (%s)</li>",
			repo.Server,
			repo.FullName,
			action,
		)
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(
		w,
		fmt.Sprintf(
			"<html><h3>Repositories:</h3><ul>%s</ul></html>",
			repositoriesHtml,
		),
	)
}
