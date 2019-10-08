package multilogger

import (
	"fmt"
	"os"

	"github.com/acarl005/stripansi"
	"github.com/sirupsen/logrus"
)

// FileHook represents a logger that will send logs (of all levels) to a file
type FileHook struct {
	*GenericHook
	Filename string
	file     *os.File
}

// NewFileHook creates a FileHook instance
func NewFileHook(fileName string, level logrus.Level, formatter logrus.Formatter) *FileHook {
	return &FileHook{
		GenericHook: &GenericHook{
			Formatter:    formatter,
			MinimumLevel: level,
		},
		Filename: fileName,
	}
}

// Fire writes logs to the configured file
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	formatted, err := hook.formatEntry(entry)
	if len(formatted) == 0 {
		return err
	}

	if hook.file == nil {
		logFileExists := false
		if _, err := os.Stat(hook.Filename); err == nil {
			logFileExists = true
		}
		if hook.file, err = os.OpenFile(hook.Filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666); err != nil {
			return fmt.Errorf("Unable to open log file %s: %v", hook.Filename, err)
		}
		if logFileExists {
			// Add a bit of whitespace before logging
			hook.file.Write([]byte("\n\n"))
		}
		hook.file.Write([]byte("### Opening log file ###\n\n"))
	}

	if _, err = hook.file.WriteString(stripansi.Strip(string(formatted))); err != nil {
		return fmt.Errorf("Unable to print logs to file: %v", err)
	}

	return nil
}
