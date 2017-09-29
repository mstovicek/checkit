package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mstovicek/checkit/document_store"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/repository_api"
	"html"
	"io"
	"net/http"
)

type inspectionDetailHandler struct {
	logger        logger.Log
	documentStore document_store.Store
	statusMap     map[int]string
}

func newInspectionDetailHandler(
	log logger.Log,
	documentStore document_store.Store,
) http.Handler {
	return &inspectionDetailHandler{
		logger:        log,
		documentStore: documentStore,
		statusMap: map[int]string{
			repository_api.CommitStatusPending: "Pending",
			repository_api.CommitStatusSuccess: "Success",
			repository_api.CommitStatusFailure: "Failure",
			repository_api.CommitStatusError:   "Error",
		},
	}
}

func (h *inspectionDetailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, exists := vars["uuid"]
	if !exists {
		writeError(w, http.StatusBadRequest, "inspection ID is not set")
		return
	}

	document := document_store.InspectionResultDocument{}
	err := h.documentStore.FindOne(document_store.QueryField{"uuid": uuid}, &document)
	if err != nil {
		writeError(w, http.StatusNotFound, "cannot find inspection details")
		return
	}

	jsonDoc, _ := json.MarshalIndent(document, "", "  ")
	inspectionHtml := "detail of " + uuid + "<br><pre>" + string(jsonDoc) + "</pre>"

	filesHtml := ""
	for _, file := range document.FixedFiles {
		filesHtml += "<b>" + file.Name + "</b><br><pre>" + html.EscapeString(file.Diff) + "</pre><br>"
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(
		w,
		fmt.Sprintf(
			"<html><h3>Inspection - %s:</h3>%s<br>%s</html>",
			h.statusMap[document.Status],
			inspectionHtml,
			filesHtml,
		),
	)
}
