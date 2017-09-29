package api

import (
	"encoding/json"
	"net/http"
)

type errorsEnvelope struct {
	Errors []string `json:"errors"`
}

func writeError(w http.ResponseWriter, code int, message string) error {
	w.WriteHeader(code)

	payload := errorsEnvelope{
		Errors: []string{message},
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(payload)
}
