package document_store

import (
	"errors"
	"github.com/mstovicek/checkit/logger"
	"gopkg.in/mgo.v2"
)

type mongoDB struct {
	logger         logger.Log
	databaseName   string
	collectionName string
	session        *mgo.Session
}

func NewMongoDB(log logger.Log, mongoDBURL string, databaseName string, collectionName string) (Store, error) {
	mgoSession, err := mgo.Dial(mongoDBURL)
	mgoSession.New()
	if err != nil {
		log.Error(logger.Fields{
			"mongoDBURL":     mongoDBURL,
			"databaseName":   databaseName,
			"collectionName": collectionName,
		}, "cannot connect to MongoDB")
		return nil, errors.New("cannot dial MongoDB")
	}

	mgoSession.SetMode(mgo.Monotonic, true)

	return &mongoDB{
		logger:         log,
		databaseName:   databaseName,
		collectionName: collectionName,
		session:        mgoSession,
	}, nil
}

func (mongoDB *mongoDB) EnsureIndex(key string) error {
	session, collection := mongoDB.copySessionWithCollection()
	defer session.Close()

	index := mgo.Index{
		Key:        []string{key},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return collection.EnsureIndex(index)
}

func (mongoDB *mongoDB) Close() {
	mongoDB.session.Close()
}

func (mongoDB *mongoDB) Insert(documents ...interface{}) error {
	session, collection := mongoDB.copySessionWithCollection()
	defer session.Close()

	return collection.Insert(documents...)
}

func (mongoDB *mongoDB) FindOne(queryFields QueryField, output interface{}) error {
	session, collection := mongoDB.copySessionWithCollection()
	defer session.Close()

	err := collection.Find(queryFields).One(output)

	mongoDB.logger.Debug(logger.Fields{
		"output":      output,
		"queryFields": queryFields,
	}, "FindOne")

	return err
}

func (mongoDB *mongoDB) copySessionWithCollection() (*mgo.Session, *mgo.Collection) {
	session := mongoDB.session.Copy()
	collection := session.DB(mongoDB.databaseName).C(mongoDB.collectionName)
	return session, collection
}
