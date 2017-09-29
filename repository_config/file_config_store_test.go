package repository_config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigFileName(t *testing.T) {
	s, _ := NewFileConfigStore("./base-path/", nil)
	store := s.(*fileConfigStore)

	path := store.getConfigFilename("server", "repository/name")
	expectedPath := "./base-path/server--repository--name.json"

	assert.Equal(t, expectedPath, path)
}

func TestGetConfigFileNameWithDashes(t *testing.T) {
	s, _ := NewFileConfigStore("./base-path/", nil)
	store := s.(*fileConfigStore)

	path := store.getConfigFilename("domain/ser-ver", "repository/na-me")
	expectedPath := "./base-path/domain--ser-ver--repository--na-me.json"

	assert.Equal(t, expectedPath, path)
}

func TestGetConfigFileNameWithSpecialChars(t *testing.T) {
	s, _ := NewFileConfigStore("~/base-path/", nil)
	store := s.(*fileConfigStore)

	path := store.getConfigFilename("ser%ve$r", "re~po$sitory/%name")
	expectedPath := "~/base-path/server--repository--name.json"

	assert.Equal(t, expectedPath, path)
}
