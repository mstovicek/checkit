package code_fixer

import (
	"encoding/json"
	"github.com/mstovicek/checkit/logger"
	"github.com/mstovicek/checkit/repository_api"
	"os/exec"
)

const (
	phpFixerExitStatusOK                      = 0  // OK.
	phpFixerExitStatusError                   = 1  // General error (or PHP minimal requirement not matched).
	phpFixerExitStatusWrongSyntax             = 4  // Some files have invalid syntax (only in dry-run mode).
	phpFixerExitStatusHasFixes                = 8  // Some files need fixing (only in dry-run mode).
	phpFixerExitStatusErrorConfiguration      = 16 // Configuration error of the application.
	phpFixerExitStatusErrorFixerConfiguration = 32 // Configuration error of a Fixer.
	phpFixerExitStatusException               = 64 // Exception raised within the application.
)

type PhpConfiguration struct{}

type phpFixer struct {
	log         logger.Log
	config      PhpConfiguration
	statusMap   map[int]int
	commandName string
}

func NewPhp(log logger.Log, config PhpConfiguration, commandName string) Fixer {
	return &phpFixer{
		log:    log,
		config: config,
		statusMap: map[int]int{
			phpFixerExitStatusOK:                      repository_api.CommitStatusSuccess,
			phpFixerExitStatusError:                   repository_api.CommitStatusError,
			phpFixerExitStatusWrongSyntax:             repository_api.CommitStatusFailure,
			phpFixerExitStatusHasFixes:                repository_api.CommitStatusFailure,
			phpFixerExitStatusErrorConfiguration:      repository_api.CommitStatusError,
			phpFixerExitStatusErrorFixerConfiguration: repository_api.CommitStatusError,
			phpFixerExitStatusException:               repository_api.CommitStatusError,
		},
		commandName: commandName,
	}
}

func (fixer *phpFixer) Run(directory string, files []string) (*Result, error) {
	command := exec.Command(
		fixer.commandName,
		"fix",
		"--diff",
		"--format=json",
		"--dry-run",
		"--using-cache=no",
		"--config=.php_cs",
	)
	command.Args = append(command.Args, prefixFilesWithDirectory(directory, files)...)

	stdout, _, exitCode := runCommand(fixer.log, command)

	var result Result
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		fixer.log.Error(logger.Fields{
			"error": err.Error(),
		}, "Cannot unmarshal result in php fixer")
		return nil, err
	}

	result.Status = fixer.statusMap[exitCode]

	return &result, nil
}
