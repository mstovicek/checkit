package api

import (
	"github.com/gorilla/mux"
	"github.com/mstovicek/checkit/logger"
	"log"
	"net/http"
)

type Server interface {
	AddHandlers(h []Handler)
	Run()
}

type Handler struct {
	path    string
	handler http.Handler
	methods []string
}

type server struct {
	listenAddress string
	log           logger.Log
	handlers      []Handler
}

func NewServer(listenAddress string, log logger.Log) Server {
	return &server{
		listenAddress: listenAddress,
		log:           log,
		handlers:      []Handler{},
	}
}

func (s *server) AddHandlers(handlers []Handler) {
	s.handlers = append(s.handlers, handlers...)
}

func (s *server) Run() {
	router := mux.NewRouter()

	for _, h := range s.handlers {
		s.log.Info(logger.Fields{
			"path":    h.path,
			"methods": h.methods,
		}, "Adding handler")

		if len(h.methods) == 0 {
			router.Handle(h.path, h.handler)
		} else {
			router.Handle(h.path, h.handler).Methods(h.methods...)
		}
	}

	handler := newRecoveryMiddleware(
		s.log,
		newLogResponseTimeMiddleware(
			s.log,
			router,
		),
	)

	s.log.Info(logger.Fields{
		"listen": s.listenAddress,
	}, "Listening on the address")

	log.Fatal(http.ListenAndServe(s.listenAddress, handler))
}
