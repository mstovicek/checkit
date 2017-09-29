package logger

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	projectDir = "github.com/mstovicek/checkit"
)

type Fields map[string]interface{}

type Log interface {
	Debug(fields Fields, msg string)
	Debugf(format string, args ...interface{})
	Info(fields Fields, msg string)
	Warn(fields Fields, msg string)
	Error(fields Fields, msg string)
	Errorf(format string, args ...interface{})
	SetDebug(isDebug bool) Log
	SetComponent(component string) Log
}

func getFilePath() string {
	_, file, line, ok := runtime.Caller(4)

	if ok {
		slash := strings.LastIndex(file, projectDir)
		if slash >= 0 {
			file = file[slash+len(projectDir)+1:]
		}
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}
