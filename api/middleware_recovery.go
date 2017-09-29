package api

import (
	"github.com/mstovicek/checkit/logger"
	"net/http"
)

type recoveryMiddleware struct {
	logger  logger.Log
	handler http.Handler
}

func newRecoveryMiddleware(l logger.Log, h http.Handler) http.Handler {
	return &recoveryMiddleware{
		logger:  l,
		handler: h,
	}
}

func (m *recoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// recover allows you to continue execution in case of panic
		if err := recover(); err != nil {
			writeError(w, http.StatusInternalServerError, "Internal Server Error")

			m.logger.Error(logger.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"error":  err,
			}, "Internal Server Error")
		}
	}()

	m.handler.ServeHTTP(w, r)
}
