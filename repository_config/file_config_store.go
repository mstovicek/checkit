package repository_config

import (
	"fmt"
	"github.com/mstovicek/checkit/logger"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type fileConfigStore struct {
	logger     logger.Log
	basePath   string
	pathRegexp *regexp.Regexp
}

func NewFileConfigStore(configBasePath string, log logger.Log) (ConfigStore, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9-/]+")
	if err != nil {
		log.Error(logger.Fields{
			"error": err.Error(),
		}, "Cannot compile regexp")
		return nil, err
	}

	return &fileConfigStore{
		logger:     log,
		basePath:   configBasePath,
		pathRegexp: reg,
	}, nil
}

func (store *fileConfigStore) HasConfig(server string, repositoryName string) bool {
	filename := store.getConfigFilename(server, repositoryName)
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func (store *fileConfigStore) LoadConfig(server string, repositoryName string) (Config, error) {
	filename := store.getConfigFilename(server, repositoryName)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		store.logger.Error(logger.Fields{
			"filename": filename,
			"error":    err.Error(),
		}, "Cannot read repository config file")
		return Config{
			Server:         server,
			RepositoryName: repositoryName,
		}, err
	}

	config, err := bytesToConfig(b)
	if err != nil {
		store.logger.Error(logger.Fields{
			"filename": filename,
			"error":    err.Error(),
		}, "Cannot unserialize repository config file")
		return Config{
			Server:         server,
			RepositoryName: repositoryName,
		}, err
	}
	return config, nil
}

func (store *fileConfigStore) SaveConfig(config Config) error {
	filename := store.getConfigFilename(config.Server, config.RepositoryName)

	configBytes, err := configToBytes(config)
	if err != nil {
		store.logger.Error(logger.Fields{
			"filename":       filename,
			"server":         config.Server,
			"repositoryName": config.RepositoryName,
			"error":          err.Error(),
		}, "Cannot serialize config")
		return err
	}

	err = ioutil.WriteFile(
		filename,
		configBytes,
		0644,
	)
	if err != nil {
		store.logger.Error(logger.Fields{
			"filename":       filename,
			"server":         config.Server,
			"repositoryName": config.RepositoryName,
			"error":          err.Error(),
		}, "Cannot save config file")
		return err
	}
	return nil
}

func (store *fileConfigStore) getConfigFilename(server string, repositoryName string) string {
	return fmt.Sprintf(
		"%s%s--%s.json",
		store.basePath,
		strings.Replace(store.removeSpecialCharacters(server), "/", "--", -1),
		strings.Replace(store.removeSpecialCharacters(repositoryName), "/", "--", -1),
	)
}

func (store *fileConfigStore) removeSpecialCharacters(str string) string {
	return store.pathRegexp.ReplaceAllString(str, "")
}
