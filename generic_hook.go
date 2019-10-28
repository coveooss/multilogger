package multilogger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type genericHookI interface {
	SetFormatter(logrus.Formatter)
	Formatter() logrus.Formatter
}

type genericHook struct {
	formatter logrus.Formatter
}

func (hook *genericHook) formatEntry(entry *logrus.Entry) ([]byte, error) {
	if hook.formatter == nil {
		hook.formatter = NewFormatter(true, os.Getenv(FormatFileEnvVar), os.Getenv(FormatEnvVar), DefaultFileFormat)
	}
	formatted, err := hook.formatter.Format(entry)
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to format the given log entry: %w", err)
	}
	return formatted, nil
}

func (hook *genericHook) Levels() []logrus.Level                  { return nil }
func (hook *genericHook) Fire(entry *logrus.Entry) error          { return fmt.Errorf("Not implemented") }
func (hook *genericHook) SetFormatter(formatter logrus.Formatter) { hook.formatter = formatter }
func (hook *genericHook) Formatter() logrus.Formatter             { return hook.formatter }
