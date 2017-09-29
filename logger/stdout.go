package logger

import (
	"github.com/sirupsen/logrus"
)

type logStdout struct {
	logrus    *logrus.Logger
	component string
}

func NewStdout() Log {
	return &logStdout{
		logrus: logrus.New(),
	}
}

func (l *logStdout) Debug(fields Fields, msg string) {
	l.newEntry(fields).Debug(msg)
}

func (l *logStdout) Debugf(format string, args ...interface{}) {
	l.newEntryEmpty().Debugf(format, args...)
}

func (l *logStdout) Info(fields Fields, msg string) {
	l.newEntry(fields).Info(msg)
}

func (l *logStdout) Warn(fields Fields, msg string) {
	l.newEntry(fields).Warn(msg)
}

func (l *logStdout) Error(fields Fields, msg string) {
	l.newEntry(fields).Error(msg)
}

func (l *logStdout) Errorf(format string, args ...interface{}) {
	l.newEntryEmpty().Errorf(format, args...)
}

func (l *logStdout) SetDebug(isDebug bool) Log {
	if isDebug {
		l.logrus.Level = logrus.DebugLevel
	} else {
		l.logrus.Level = logrus.InfoLevel
	}

	return l
}

func (l *logStdout) SetComponent(component string) Log {
	l.component = component

	return l
}

func (l *logStdout) newEntry(fields Fields) *logrus.Entry {
	entry := l.newEntryEmpty()

	for key, value := range fields {
		entry = entry.WithField(key, value)
	}

	return entry
}

func (l *logStdout) newEntryEmpty() *logrus.Entry {
	entry := logrus.NewEntry(l.logrus)
	entry = entry.WithField("file", getFilePath())

	if l.component != "" {
		entry = entry.WithField("component", l.component)
	}

	return entry
}
