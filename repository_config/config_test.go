package repository_config

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func TestBytesToConfig(t *testing.T) {
	b := []byte(`{
	"server": "server",
	"name": "repoName",
	"url": "http://url",
	"clone_url": "http://a:b@url.git",
	"email": "email@foo.bar"
}`)
	actualConfig, _ := bytesToConfig(b)

	expectedConfig := Config{
		Server:         "server",
		RepositoryName: "repoName",
		Url:            "http://url",
		CloneUrl:       "http://a:b@url.git",
		Email:          "email@foo.bar",
	}

	assert.Equal(t, expectedConfig, actualConfig)
}

func TestBytesToConfigEmpty(t *testing.T) {
	b := []byte("{}")
	actualConfig, _ := bytesToConfig(b)

	expectedConfig := Config{}

	assert.Equal(t, expectedConfig, actualConfig)
}

func TestBytesToConfigWithToken(t *testing.T) {
	b := []byte(`{
	"server": "server",
	"name": "repoName",
	"url": "http://url",
	"clone_url": "http://a:b@url.git",
	"email": "email@foo.bar",
	"oauth_token": {
		"access_token": "accessToken"
	}
}`)
	actualConfig, _ := bytesToConfig(b)

	expectedConfig := Config{
		Server:         "server",
		RepositoryName: "repoName",
		Url:            "http://url",
		CloneUrl:       "http://a:b@url.git",
		Email:          "email@foo.bar",
		OAuthToken: oauth2.Token{
			AccessToken: "accessToken",
		},
	}

	assert.Equal(t, expectedConfig, actualConfig)
}

func TestConfigToBytes(t *testing.T) {
	c := Config{
		Server:         "server",
		RepositoryName: "repoName",
		Url:            "http://url",
		CloneUrl:       "http://a:b@url.git",
		Email:          "email@foo.bar",
	}
	actualBytes, _ := configToBytes(c)

	expectedBytes := []byte(`{
	"server": "server",
	"name": "repoName",
	"url": "http://url",
	"clone_url": "http://a:b@url.git",
	"email": "email@foo.bar",
	"oauth_token": {
		"access_token": "",
		"expiry": "0001-01-01T00:00:00Z"
	}
}`)

	assert.Equal(t, string(expectedBytes), string(actualBytes))
	assert.Equal(t, expectedBytes, actualBytes)
}

func TestConfigToBytesWithToken(t *testing.T) {
	c := Config{
		Server:         "server",
		RepositoryName: "repoName",
		Url:            "http://url",
		CloneUrl:       "http://a:b@url.git",
		Email:          "email@foo.bar",
		OAuthToken: oauth2.Token{
			AccessToken: "accessToken",
		},
	}
	actualBytes, _ := configToBytes(c)

	expectedBytes := []byte(`{
	"server": "server",
	"name": "repoName",
	"url": "http://url",
	"clone_url": "http://a:b@url.git",
	"email": "email@foo.bar",
	"oauth_token": {
		"access_token": "accessToken",
		"expiry": "0001-01-01T00:00:00Z"
	}
}`)

	assert.Equal(t, string(expectedBytes), string(actualBytes))
	assert.Equal(t, expectedBytes, actualBytes)
}
