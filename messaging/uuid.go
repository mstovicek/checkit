package messaging

import (
	"github.com/satori/go.uuid"
)

type UUIDGenerator interface {
	Generate() string
}

func NewUUIDGenerator() UUIDGenerator {
	return &uuidGenerator{}
}

type uuidGenerator struct{}

func (g *uuidGenerator) Generate() string {
	return uuid.NewV4().String()
}
