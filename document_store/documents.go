package document_store

import "time"

type InspectionResultDocument struct {
	Uuid           string    `bson:"uuid"`
	Time           time.Time `bson:"time"`
	Server         string    `bson:"server"`
	RepositoryName string    `bson:"repository_name"`
	CommitHash     string    `bson:"commit_hash"`
	Status         int       `bson:"status"`
	FixedFiles     []struct {
		Name string // `bson:"name"`
		Diff string // `bson:"diff"`
	} `bson:"fixed_files"`
}
