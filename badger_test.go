package benchmarks_test

import (
	"fmt"
	"log/slog"
	"strings"
)

type badgerLogger struct {
	log *slog.Logger
}

func newBadgerLogger(log *slog.Logger) *badgerLogger {
	return &badgerLogger{log: log.With("code.namespace", "badger-lib")}
}

func (l *badgerLogger) Errorf(f string, v ...interface{}) {
	l.log.Error(fmt.Sprintf(strings.Trim(f, "\n"), v...))
}

func (l *badgerLogger) Warningf(f string, v ...interface{}) {
	l.log.Warn(fmt.Sprintf(strings.Trim(f, "\n"), v...))
}

func (l *badgerLogger) Infof(f string, v ...interface{}) {
	l.log.Debug(fmt.Sprintf(strings.Trim(f, "\n"), v...))
}

func (l *badgerLogger) Debugf(f string, v ...interface{}) {
	l.log.Debug(fmt.Sprintf(strings.Trim(f, "\n"), v...))
}
