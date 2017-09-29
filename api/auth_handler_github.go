package api

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"github.com/mstovicek/checkit/repository_api"
	"net/http"
)

type authGithubHandler struct {
	logger           logger.Log
	sessionStore     sessions.Store
	sessionStoreName string
}

func newAuthGithubHandler(logger logger.Log, sessionStore sessions.Store, sessionStoreName string) http.Handler {
	return &authGithubHandler{
		logger:           logger,
		sessionStore:     sessionStore,
		sessionStoreName: sessionStoreName,
	}
}

func (h *authGithubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionStoreName)
	if err != nil {
		h.logger.Error(logger.Fields{
			"sessionStoreName": h.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot get session store")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	sessionState := session.Values[sessionKeyState]
	queryState := r.URL.Query().Get(queryKeyState)
	if sessionState != queryState {
		h.logger.Error(logger.Fields{
			"sessionState": sessionState,
			"queryState":   queryState,
		}, "Github OAuth state mismatch")
		writeError(w, http.StatusBadRequest, "State mismatch")
		return
	}

	code := r.URL.Query().Get("code")
	auth, err := oauth.NewOAuth(oauth.GithubServer, h.logger)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Unsupported OAuth server")
		return
	}

	token, err := auth.ExchangeAccessToken(code)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot exchange access token")
		return
	}

	repositoryApi := repository_api.NewGithubUserAPI(
		auth.Config.Client(context.Background(), token),
		h.logger,
	)
	user, err := repositoryApi.GetUser()
	if err != nil {
		h.logger.Error(logger.Fields{
			"server": oauth.GithubServer,
			"token":  token,
			"error":  err.Error(),
		}, "Cannot find any email")
		writeError(w, http.StatusInternalServerError, "Cannot find any email")
		return
	}

	session.Values[sessionKeyEmail] = user.Email
	session.Values[sessionKeyOAuthServer] = oauth.GithubServer
	session.Values[sessionKeyAccessToken], err = oauth.EncodeToken(token)
	if err != nil {
		h.logger.Error(logger.Fields{
			"server": oauth.GithubServer,
			"email":  user.Email,
			"token":  token,
			"error":  err.Error(),
		}, "Cannot encode access token")
		writeError(w, http.StatusInternalServerError, "Cannot encode access token")
		return
	}

	err = session.Save(r, w)
	if err != nil {
		h.logger.Error(logger.Fields{
			"sessionStoreName": h.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot save session")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	http.Redirect(w, r, "/auth/", http.StatusFound)
}
