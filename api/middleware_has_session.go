package api

import (
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"net/http"
)

type hasSessionMiddleware struct {
	logger           logger.Log
	sessionStore     sessions.Store
	sessionStoreName string
	handler          http.Handler
}

func newHasSessionMiddleware(
	log logger.Log,
	sessionStore sessions.Store,
	sessionStoreName string,
	h http.Handler,
) http.Handler {
	return &hasSessionMiddleware{
		logger:           log,
		sessionStore:     sessionStore,
		sessionStoreName: sessionStoreName,
		handler:          h,
	}
}

func (m *hasSessionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := m.sessionStore.Get(r, m.sessionStoreName)
	if err != nil {
		m.logger.Error(logger.Fields{
			"sessionStoreName": m.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot get session store")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	if session.Values[sessionKeyEmail] == nil || session.Values[sessionKeyAccessToken] == nil || session.Values[sessionKeyOAuthServer] == nil {
		m.logger.Info(logger.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}, "Unathorized")
		writeError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	m.handler.ServeHTTP(w, r)
}
