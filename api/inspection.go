package api

import (
	"github.com/gorilla/sessions"
	"github.com/mstovicek/checkit/document_store"
	"github.com/mstovicek/checkit/logger"
	"net/http"
)

func GetInspectionHandlers(
	baseUrl string,
	log logger.Log,
	sessionAuthenticationKey string,
	sessionEncryptionKey string,
	sessionStoreName string,
	mongoDBURL string,
	mongoDBDatabaseName string,
	mongoDBCollectionName string,
) ([]Handler, error) {
	sessionStore := sessions.NewCookieStore(
		[]byte(sessionAuthenticationKey),
		[]byte(sessionEncryptionKey),
	)

	documentStore, err := document_store.NewMongoDB(log, mongoDBURL, mongoDBDatabaseName, mongoDBCollectionName)
	if err != nil {
		return nil, err
	}
	err = documentStore.EnsureIndex("uuid")
	if err != nil {
		return nil, err
	}

	return []Handler{
		{
			path: baseUrl + "{uuid}/",
			handler: newHasSessionMiddleware(
				log,
				sessionStore,
				sessionStoreName,
				newInspectionDetailHandler(
					log,
					documentStore,
				),
			),
			methods: []string{http.MethodGet},
		},
	}, nil
}
