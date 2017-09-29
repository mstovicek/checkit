package git

import (
	"fmt"
	"github.com/mstovicek/checkit/logger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type filesystemClient struct {
	basePath string
	log      logger.Log
}

type repositoryInformation struct {
	cloneUrl       string
	server         string
	repositoryName string
	commitHash     string
	directory      string
}

func NewFilesystemClient(basePath string, log logger.Log) Client {
	return &filesystemClient{
		basePath: basePath,
		log:      log,
	}
}

func (c *filesystemClient) Checkout(cloneUrl string, server string, repositoryName string, commitHash string) (string, error) {
	directory := c.getConfigFilename(server, repositoryName)

	repoInfo := repositoryInformation{
		cloneUrl:       cloneUrl,
		server:         server,
		repositoryName: repositoryName,
		commitHash:     commitHash,
		directory:      directory,
	}

	repository, err := c.clone(repoInfo)
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		c.log.Error(logger.Fields{
			"error": err.Error(),
		}, "cannot clone repo")
		return "", err
	}

	if err == git.ErrRepositoryAlreadyExists {
		repository, err = c.open(repoInfo)
		if err != nil {
			c.log.Error(logger.Fields{
				"error": err.Error(),
			}, "cannot open repo")
			return "", err
		}
	}

	if err := c.fetch(repository, repoInfo); err != nil {
		c.log.Error(logger.Fields{
			"error": err.Error(),
		}, "cannot fetch repo")
		return "", err
	}

	if err := c.resetHead(repository, repoInfo); err != nil {
		c.log.Error(logger.Fields{
			"error": err.Error(),
		}, "cannot reset HEAD")
		return "", err
	}

	if err := c.checkout(repository, commitHash, repoInfo); err != nil {
		c.log.Error(logger.Fields{
			"error": err.Error(),
		}, "cannot checkout repo")
		return "", err
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := repository.Head()
	if err != nil {
		c.log.Error(logger.Fields{
			"error": err.Error(),
		}, "cannot get HEAD")
	}

	c.log.Debug(logger.Fields{
		"server":         server,
		"repositoryName": repositoryName,
		"commitHash":     ref.Hash().String(),
	}, "git checkout done")

	return directory, err
}

func (c *filesystemClient) clone(repoInfo repositoryInformation) (*git.Repository, error) {
	c.log.Debug(c.repoInfoToLoggerFields(repoInfo), "cloning repository")
	return git.PlainClone(
		repoInfo.directory,
		false,
		&git.CloneOptions{
			URL:               repoInfo.cloneUrl,
			RecurseSubmodules: 0,
			Depth:             1,
		},
	)
}

func (c *filesystemClient) open(repoInfo repositoryInformation) (*git.Repository, error) {
	c.log.Debug(c.repoInfoToLoggerFields(repoInfo), "opening repository")
	return git.PlainOpen(repoInfo.directory)
}

func (c *filesystemClient) fetch(repository *git.Repository, repoInfo repositoryInformation) error {
	c.log.Debug(c.repoInfoToLoggerFields(repoInfo), "fetching repository")
	err := repository.Fetch(&git.FetchOptions{})
	if err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func (c *filesystemClient) resetHead(repository *git.Repository, repoInfo repositoryInformation) error {
	c.log.Debug(c.repoInfoToLoggerFields(repoInfo), "Reseting hard -- not implemented in the library :/")

	return nil

	//w, err := repository.Worktree()
	//if err != nil {
	//	c.log.Error(logger.Fields{
	//		"error": err.Error(),
	//	}, "cannot get working tree")
	//	return err
	//}
	//
	//err = w.Reset(&git.ResetOptions{Mode: git.HardReset})
	//if err != nil {
	//	c.log.Error(logger.Fields{
	//		"error": err.Error(),
	//	}, "cannot reset HEAD")
	//	return err
	//}
	//return nil
}

func (c *filesystemClient) checkout(repository *git.Repository, commitHash string, repoInfo repositoryInformation) error {
	c.log.Debug(c.repoInfoToLoggerFields(repoInfo), "Checking out")
	c.log.Debug(logger.Fields{"commitHash": commitHash}, "Checking out commit")

	w, err := repository.Worktree()
	if err != nil {
		c.log.Error(logger.Fields{
			"error": err.Error(),
		}, "cannot get working tree")
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})
	if err != nil {
		c.log.Error(logger.Fields{
			"error":      err.Error(),
			"commitHash": commitHash,
		}, "cannot checkout commit")
		return err
	}

	return nil
}

func (c *filesystemClient) getConfigFilename(server string, repositoryName string) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		c.basePath,
		server,
		repositoryName,
	)
}

func (c *filesystemClient) repoInfoToLoggerFields(repoInfo repositoryInformation) logger.Fields {
	return logger.Fields{
		"cloneUrl":       repoInfo.cloneUrl,
		"server":         repoInfo.server,
		"repositoryName": repoInfo.repositoryName,
		"commitHash":     repoInfo.commitHash,
		"directory":      repoInfo.directory,
	}
}
