package api

import (
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/logger"
	"net/http"
)

type authLogoutHandler struct {
	logger           logger.Log
	sessionStore     sessions.Store
	sessionStoreName string
}

func newAuthLogoutHandler(logger logger.Log, sessionStore sessions.Store, sessionStoreName string) http.Handler {
	return &authLogoutHandler{
		logger:           logger,
		sessionStore:     sessionStore,
		sessionStoreName: sessionStoreName,
	}
}

func (h *authLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := h.sessionStore.Get(r, h.sessionStoreName)
	if err != nil {
		h.logger.Error(logger.Fields{
			"sessionStoreName": h.sessionStoreName,
			"error":            err.Error(),
		}, "Cannot get session store")
		writeError(w, http.StatusInternalServerError, "Cannot process a request")
		return
	}

	session.Options.MaxAge = -1

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
