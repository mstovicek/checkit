package logger

import (
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"net"
)

type logLogstash struct {
	logrus    *logrus.Logger
	component string
}

func NewLogstash(logstashTcpAddress string) Log {
	l := logrus.New()

	conn, err := net.Dial("tcp", logstashTcpAddress)
	if err != nil {
		panic(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{}))

	if err != nil {
		panic(err)
	}
	l.Hooks.Add(hook)

	return &logLogstash{
		logrus: l,
	}
}

func (l *logLogstash) Debug(fields Fields, msg string) {
	l.newEntry(fields).Debug(msg)
}

func (l *logLogstash) Debugf(format string, args ...interface{}) {
	l.newEntryEmpty().Debugf(format, args...)
}

func (l *logLogstash) Info(fields Fields, msg string) {
	l.newEntry(fields).Info(msg)
}

func (l *logLogstash) Warn(fields Fields, msg string) {
	l.newEntry(fields).Warn(msg)
}

func (l *logLogstash) Error(fields Fields, msg string) {
	l.newEntry(fields).Error(msg)
}

func (l *logLogstash) Errorf(format string, args ...interface{}) {
	l.newEntryEmpty().Errorf(format, args...)
}

func (l *logLogstash) SetDebug(isDebug bool) Log {
	if isDebug {
		l.logrus.Level = logrus.DebugLevel
	} else {
		l.logrus.Level = logrus.InfoLevel
	}

	return l
}

func (l *logLogstash) SetComponent(component string) Log {
	l.component = component

	return l
}

func (l *logLogstash) newEntry(fields Fields) *logrus.Entry {
	entry := l.newEntryEmpty()

	for key, value := range fields {
		entry = entry.WithField(key, value)
	}

	return entry
}

func (l *logLogstash) newEntryEmpty() *logrus.Entry {
	entry := logrus.NewEntry(l.logrus)
	entry = entry.WithField("file", getFilePath())

	if l.component != "" {
		entry = entry.WithField("component", l.component)
	}

	return entry
}
