package multilogger

import (
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

func (hook *consoleHook) clone() logrus.Hook {
	// Duplicate the console hook to ensure that the copy
	// has its own attributes when the object is copied.
	return &consoleHook{
		genericHook: hook.genericHook,
		out:         hook.out,
		log:         hook.log,
	}
}

func (hook *consoleHook) Fire(entry *logrus.Entry) (err error) {
	return hook.fire(entry, func() error {
		const name = "ConsoleHook"
		if entry.Level == outputLevel {
			return hook.printf(name, hook.out, entry.Message)
		}
		var formatted string
		if formatted, err = hook.formatEntry(name, entry); err != nil {
			return err
		}
		return hook.printf(name, hook.log, string(formatted))
	})
}

func (hook *consoleHook) SetOut(out io.Writer)    { hook.log = out }
func (hook *consoleHook) SetStdout(out io.Writer) { hook.out = out }
