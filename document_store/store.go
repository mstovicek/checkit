package document_store

type Store interface {
	EnsureIndex(key string) error
	Close()
	Insert(document ...interface{}) error
	FindOne(queryFields QueryField, output interface{}) error
}

type QueryField map[string]interface{}
