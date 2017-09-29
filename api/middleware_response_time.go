package api

import (
	"github.com/mstovicek/checkit/logger"
	"net/http"
	"time"
)

type logResponseTimeMiddleware struct {
	logger  logger.Log
	handler http.Handler
}

func newLogResponseTimeMiddleware(l logger.Log, h http.Handler) http.Handler {
	return &logResponseTimeMiddleware{
		logger:  l,
		handler: h,
	}
}

func (m *logResponseTimeMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()

	m.handler.ServeHTTP(w, r)

	m.logger.Info(logger.Fields{
		"method":        r.Method,
		"path":          r.URL.Path,
		"response_time": time.Since(started).Nanoseconds(),
	}, "Request has been processed")
}
