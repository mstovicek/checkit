package messaging

import (
	"encoding/json"
	"github.com/mstovicek/checkit/logger"
	"time"
)

type CommitReceived struct {
	Time           time.Time `json:"time"`
	Server         string    `json:"server"`
	RepositoryName string    `json:"repository_name"`
	CommitHash     string    `json:"commit_hash"`
	Files          []string  `json:"files"`
}

type InspectionProcessed struct {
	Time           time.Time `json:"time"`
	Server         string    `json:"server"`
	RepositoryName string    `json:"repository_name"`
	CommitHash     string    `json:"commit_hash"`
	Status         int       `json:"status"`
	FixedFiles     []struct {
		Name string // `json:"name"`
		Diff string // `json:"diff"`
	} `json:"fixed_files"`
}

func Unmarshal(deliveredBody []byte, output interface{}, log logger.Log) error {
	err := json.Unmarshal(deliveredBody, output)
	if err != nil {
		log.Error(logger.Fields{"error": err.Error()}, "cannot unmarshal message body")
		return err
	}
	return nil
}
