package multilogger

import (
	"io"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultConsoleFormat is the format used by NewConsoleHook if MULTILOGGER_FORMAT is not set.
	DefaultConsoleFormat = "%module:Italic,Green,SquareBrackets,IgnoreEmpty,Space%%time% %-8level:upper,color% %message:color%"
	consoleHookName      = "console-hook"
)

// GetDefaultConsoleHook returns the default console hook.
func (logger *Logger) GetDefaultConsoleHook() *Hook { return logger.Hook(consoleHookName) }

// GetDefaultConsoleHookLevel returns the logging level associated to the default console hook.
func (logger *Logger) GetDefaultConsoleHookLevel() logrus.Level {
	return logger.GetHookLevel(consoleHookName)
}

// SetDefaultConsoleHookLevel set a new log level for the default console hook.
func (logger *Logger) SetDefaultConsoleHookLevel(level interface{}) error {
	return logger.SetHookLevel(consoleHookName, level)
}

// GetFormatter returns the formater associated to the default console hook.
func (logger *Logger) GetFormatter() (result logrus.Formatter) {
	logger.onDefaultHook(func(ch *Hook) { result = ch.GetFormatter() })
	return
}

// Formatter returns the Formatter associated to the default console hook.
func (logger *Logger) Formatter() (result *Formatter) {
	logger.onDefaultHook(func(ch *Hook) { result = ch.Formatter() })
	return
}

// SetFormatter allows setting a formatter to the default console hook if there is.
func (logger *Logger) SetFormatter(formatter logrus.Formatter) *Logger {
	return logger.onDefaultHook(func(ch *Hook) { ch.SetFormatter(formatter) })
}

// SetFormat allows setting a format string on the default console hook if there is.
func (logger *Logger) SetFormat(formats ...interface{}) *Logger {
	return logger.onDefaultHook(func(ch *Hook) { ch.SetFormat(formats...) })
}

// SetColor allows setting color mode on the default console hook if there is.
func (logger *Logger) SetColor(color bool) *Logger {
	return logger.onDefaultHook(func(ch *Hook) { ch.SetColor(color) })
}

// SetOut allows configuring the logging stream on the default console hook if there is.
func (logger *Logger) SetOut(out io.Writer) *Logger {
	return logger.onDefaultHook(func(ch *Hook) { ch.SetOut(out) })
}

// SetStdout allows configuring the output stream on the default console hook if there is.
func (logger *Logger) SetStdout(out io.Writer) *Logger {
	return logger.onDefaultHook(func(ch *Hook) { ch.SetStdout(out) })
}

// SetAllOutputs allows configuring the output and the logging stream on the default console hook if there is.
func (logger *Logger) SetAllOutputs(out io.Writer) *Logger {
	return logger.onDefaultHook(func(ch *Hook) { ch.SetOut(out).SetStdout(out) })
}

// GetDefaultInnerHook returns the inner hook actually used by the default console hook.
func (logger *Logger) GetDefaultInnerHook() (result logrus.Hook) {
	logger.onDefaultHook(func(ch *Hook) { result = ch.GetInnerHook() })
	return
}

func (logger *Logger) onDefaultHook(action func(*Hook)) *Logger {
	if ch := logger.GetDefaultConsoleHook(); ch != nil {
		action(ch)
	}
	return logger
}
