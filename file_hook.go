package multilogger

import (
	"fmt"
	"os"
	"sync"

	"github.com/acarl005/stripansi"
	"github.com/sirupsen/logrus"
)

var fileMutex sync.Mutex

// FileHook represents a logger that will send logs (of all levels) to a file
type FileHook struct {
	*GenericHook
	Filename string
	file     *os.File
}

// NewFileHook creates a FileHook instance
func NewFileHook(filename string, level logrus.Level, formatter logrus.Formatter) *FileHook {
	fileHook := &FileHook{
		GenericHook: &GenericHook{
			Formatter:    formatter,
			MinimumLevel: level,
		},
	}
	fileHook.SetFilename(filename)
	return fileHook
}

// SetFilename modifies the target file name of the hook
func (hook *FileHook) SetFilename(filename string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if hook.file != nil {
		hook.file.Close()
	}
	hook.file = nil
	hook.Filename = filename

}

// Fire writes logs to the configured file
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	formatted, err := hook.formatEntry(entry)
	if len(formatted) == 0 {
		return err
	}

	fileMutex.Lock()
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
	fileMutex.Unlock()

	if _, err = hook.file.WriteString(stripansi.Strip(string(formatted))); err != nil {
		return fmt.Errorf("Unable to print logs to file: %v", err)
	}

	return nil
}
