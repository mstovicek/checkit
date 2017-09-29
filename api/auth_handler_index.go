package api

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/oauth"
	"io"
	"net/http"
)

type authIndexHandler struct {
	logger           logger.Log
	sessionStore     sessions.Store
	sessionStoreName string
}

func newAuthIndexHandler(logger logger.Log, sessionStore sessions.Store, sessionStoreName string) http.Handler {
	return &authIndexHandler{
		logger:           logger,
		sessionStore:     sessionStore,
		sessionStoreName: sessionStoreName,
	}
}

func (h *authIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionStoreName)
	if err != nil {
		h.logger.Error(logger.Fields{
			"sessionStoreName": h.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot get session store")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	if session.Values[sessionKeyEmail] != nil && session.Values[sessionKeyAccessToken] != nil && session.Values[sessionKeyOAuthServer] != nil {
		w.WriteHeader(http.StatusOK)
		io.WriteString(
			w,
			fmt.Sprintf(
				"<html>Logged as: %s (%s)",
				session.Values[sessionKeyEmail],
				"<a href='/auth/logout/'>logout</a>",
			),
		)
	} else {
		auth, err := oauth.NewOAuth(oauth.GithubServer, h.logger)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Unsupported OAuth server")
		}

		state := auth.GetRandomState()

		session.Values[sessionKeyState] = state
		err = session.Save(r, w)
		if err != nil {
			h.logger.Error(logger.Fields{
				"sessionStoreName": h.sessionStoreName,
				"error":            err.Error(),
			}, "Cannot save session")
			writeError(w, http.StatusInternalServerError, "Cannot process a request")
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "<a href='"+auth.GetAuthUrl(state)+"'><button>Login with Github!</button></a>")
	}
}
