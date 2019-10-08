package multilogger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// ConsoleHook represents a logger that will write logs to the console
type ConsoleHook struct {
	*GenericHook
}

// NewConsoleHook creates a ConsoleHook instance
func NewConsoleHook(level logrus.Level, formatter logrus.Formatter) *ConsoleHook {
	return &ConsoleHook{
		GenericHook: &GenericHook{
			Formatter:    formatter,
			MinimumLevel: level,
		},
	}
}

// Fire writes logs to the console
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	formatted, err := hook.formatEntry(entry)
	if len(formatted) == 0 {
		return err
	}
	colorFunc := getColor(entry.Level)
	formattedString := colorFunc(string(formatted))

	if _, err = os.Stderr.WriteString(formattedString); err != nil {
		return fmt.Errorf("Unable to print logs to file: %v", err)
	}

	return nil
}
