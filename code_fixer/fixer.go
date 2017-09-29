package code_fixer

import (
	"bytes"
	"github.com/mstovicek/checkit/logger"
	"os/exec"
	"syscall"
)

type Fixer interface {
	Run(directory string, files []string) (*Result, error)
}

type Result struct {
	Status     int `json:"status"`
	FixedFiles []struct {
		Name string // `json:"name"`
		Diff string // `json:"diff"`
	} `json:"files"`
	Memory int `json:"memory"`
	Time   struct {
		Total float64 `json:"total"`
	} `json:"time"`
}

func prefixFilesWithDirectory(directory string, files []string) []string {
	for i, file := range files {
		files[i] = directory + "/" + file
	}
	return files
}

func runCommand(log logger.Log, cmd *exec.Cmd) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()

	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	log.Debug(logger.Fields{
		"commandPath": cmd.Path,
		"commandArgs": cmd.Args,
		"stdout":      stdout,
		"stderr":      stderr,
		"exitCode":    exitCode,
	}, "Processed command")

	return stdout, stderr, exitCode
}
