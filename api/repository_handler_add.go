package api

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"github.com/mstovicek/checkit/repository_api"
	"github.com/mstovicek/checkit/repository_config"
	"io"
	"net/http"
)

type repositoryGithubHandler struct {
	logger             logger.Log
	sessionStore       sessions.Store
	sessionStoreName   string
	repoConfigBasePath string
}

func newAddRepositoryHandler(logger logger.Log, sessionStore sessions.Store, sessionStoreName string, repoConfigBasePath string) http.Handler {
	return &repositoryGithubHandler{
		logger:             logger,
		sessionStore:       sessionStore,
		sessionStoreName:   sessionStoreName,
		repoConfigBasePath: repoConfigBasePath,
	}
}

func (h *repositoryGithubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionStoreName)
	if err != nil {
		h.logger.Error(logger.Fields{
			"sessionStoreName": h.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot get session store")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	repositoryName := r.URL.Query().Get("repository")
	if repositoryName == "" {
		writeError(w, http.StatusBadRequest, "Repository parameter is required")
		return
	}

	auth, err := oauth.NewOAuth(session.Values[sessionKeyOAuthServer].(string), h.logger)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Unsupported OAuth server")
		return
	}

	accessToken, err := oauth.DecodeToken(session.Values[sessionKeyAccessToken])
	if err != nil {
		h.logger.Error(logger.Fields{
			"repositoryName": repositoryName,
			"error":          err.Error(),
		}, "Cannot decode access token")
		writeError(w, http.StatusInternalServerError, "Cannot decode access token")
		return
	}

	repositoryApi, err := repository_api.NewRepositoryAPI(
		session.Values[sessionKeyOAuthServer].(string),
		auth.Config.Client(context.Background(), accessToken),
		h.logger,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Cannot connect to repository API")
		return
	}
	repo, err := repositoryApi.GetRepository(repositoryName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot read repository")
		return
	}

	if !repo.IsAdmin {
		writeError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	accessToken, err = oauth.DecodeToken(session.Values[sessionKeyAccessToken])
	if err != nil {
		h.logger.Error(logger.Fields{
			"repositoryName": repositoryName,
			"error":          err.Error(),
		}, "Cannot decode access token")
		writeError(w, http.StatusInternalServerError, "Cannot decode access token")
		return
	}

	newConfig := repository_config.Config{
		Server:         session.Values[sessionKeyOAuthServer].(string),
		RepositoryName: repositoryName,
		Url:            repo.URL,
		CloneUrl:       repo.CloneURL,
		Email:          session.Values[sessionKeyEmail].(string),
		OAuthToken:     *accessToken,
	}

	configStore, err := repository_config.NewFileConfigStore(h.repoConfigBasePath, h.logger)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot initialize repository config store")
		return
	}

	err = configStore.SaveConfig(newConfig)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot store repository config")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{\"status\": \"ok\"}")
}
