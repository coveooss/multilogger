package multilogger

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type consoleI interface {
	SetOut(io.Writer)
	SetStdout(io.Writer)
}

type consoleHook struct {
	*genericHook
	out io.Writer
	log io.Writer
}

func (hook *consoleHook) Fire(entry *logrus.Entry) error {
	if entry.Level == outputLevel {
		_, err := hook.out.Write([]byte(entry.Message))
		return err
	}
	if formatted, err := hook.formatEntry(entry); err != nil {
		return err
	} else if _, err = hook.log.Write(formatted); err != nil {
		return fmt.Errorf("Unable to fire entry: %w", err)
	}
	return nil
}

func (hook *consoleHook) SetOut(out io.Writer)    { hook.log = out }
func (hook *consoleHook) SetStdout(out io.Writer) { hook.out = out }
