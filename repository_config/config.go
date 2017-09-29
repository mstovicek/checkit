package repository_config

import (
	"encoding/json"
	"golang.org/x/oauth2"
)

type Config struct {
	Server         string       `json:"server"`
	RepositoryName string       `json:"name"`
	Url            string       `json:"url"`
	CloneUrl       string       `json:"clone_url"`
	Email          string       `json:"email"`
	OAuthToken     oauth2.Token `json:"oauth_token"`
}

func bytesToConfig(b []byte) (Config, error) {
	var conf Config
	err := json.Unmarshal(b, &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}

func configToBytes(config Config) ([]byte, error) {
	jsonConfig, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return nil, err
	}
	return jsonConfig, nil
}
