package multilogger

import (
	"fmt"
	"io"
	"os"

	"github.com/coveooss/multilogger/errors"
	"github.com/sirupsen/logrus"
)

type setLoggerI interface {
	SetLogger(*Logger)
}

type genericHookI interface {
	SetFormatter(logrus.Formatter)
	Formatter() logrus.Formatter
}

type genericHook struct {
	formatter logrus.Formatter
	logger    *Logger
}

func (hook genericHook) clone() *genericHook {
	return &hook
}

func (hook *genericHook) formatEntry(name string, entry *logrus.Entry) (string, error) {
	if hook.formatter == nil {
		hook.formatter = NewFormatter(true, os.Getenv(FormatFileEnvVar), os.Getenv(FormatEnvVar), DefaultFileFormat)
	}
	formatted, err := hook.formatter.Format(entry)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return string(formatted), nil
}

func (hook *genericHook) printf(source string, out io.Writer, format string, args ...interface{}) error {
	text := format
	if len(args) > 0 {
		text = fmt.Sprintf(format, args...)
	}
	if n, err := out.Write([]byte(text)); err != nil {
		return fmt.Errorf("%s: %w", source, err)
	} else if n != len(text) {
		return fmt.Errorf("%s: Wrong number of bytes written (%d) for %q", source, n, text)
	}
	return nil
}

func (hook *genericHook) fire(entry *logrus.Entry, f func() error) (err error) {
	defer func() {
		err = errors.Trap(err, recover())
		if hook.logger != nil && err != nil {
			// We report the error to the logger since the fire mechanism does not
			// handle errors very well
			hook.logger.AddError(err)
		}
	}()

	return f()
}

func (hook *genericHook) Levels() []logrus.Level                  { return nil }
func (hook *genericHook) Fire(entry *logrus.Entry) error          { return fmt.Errorf("Not implemented") }
func (hook *genericHook) SetFormatter(formatter logrus.Formatter) { hook.formatter = formatter }
func (hook *genericHook) Formatter() logrus.Formatter             { return hook.formatter }
func (hook *genericHook) SetLogger(l *Logger)                     { hook.logger = l }
