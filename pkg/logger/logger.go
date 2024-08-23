package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var entry *logrus.Entry

type Logger interface {
	Wrap(format string, args ...interface{}) Logger

	Info()
	Debug()
	Warn()
	Error()
}

type logger struct {
	Prefix     string
	StackTrace string
	Message    string
}

func Init(lvl string) error {
	logger := logrus.New()

	log.Println(lvl)
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	log.Println(level)

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	logrus.NewEntry(logger)

	entry = logrus.NewEntry(logger)
	return nil
}

func WithPrefix(prefix string) Logger {
	return &logger{
		Prefix: prefix,
	}
}

func (l *logger) Wrap(format string, args ...interface{}) Logger {
	l.StackTrace = l.getStackTrace()
	l.Message = fmt.Sprintf(format, args...)
	return l
}

func (l logger) Info() {
	entry.WithFields(l.extract()).Infoln(l.Message)
}

func (l logger) Debug() {
	entry.WithFields(l.extract()).Debugln(l.Message)
}

func (l logger) Warn() {
	entry.WithFields(l.extract()).Warnln(l.Message)
}

func (l logger) Error() {
	entry.WithFields(l.extract()).Errorln(l.Message)
}

func (l logger) extract() map[string]interface{} {
	return map[string]interface{}{
		"prefix": l.Prefix,
		"stack":  l.StackTrace,
	}
}

func (l logger) getStackTrace() string {
	fpcs := make([]uintptr, 32)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		log.Println("no caller")
	}

	c := runtime.FuncForPC(fpcs[0] - 1)
	if c == nil {
		log.Println("invalid caller")
		return ""
	}
	filepath, line := c.FileLine(fpcs[0] - 1)
	p, _ := os.Getwd()
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(filepath, p), line)
}
