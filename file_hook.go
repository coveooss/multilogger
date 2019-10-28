package multilogger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type fileHook struct {
	*genericHook
	filename  string
	file      *os.File
	addHeader bool
}

func (hook *fileHook) Fire(entry *logrus.Entry) (err error) {
	var output []byte
	if entry.Level == outputLevel {
		output = []byte(entry.Message)
	} else {
		if output, err = hook.formatEntry(entry); err != nil {
			return err
		}
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()
	if hook.file == nil {
		logFileExists := false
		if _, err := os.Stat(hook.filename); err == nil {
			logFileExists = true
		}
		if hook.file, err = os.OpenFile(hook.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666); err != nil {
			return fmt.Errorf("Unable to open log file %s: %w", hook.filename, err)
		}
		if hook.addHeader {
			if logFileExists {
				// Add a bit of whitespace before logging
				hook.file.Write([]byte("\n"))
			}
			hook.file.Write([]byte(fmt.Sprintf("# Opening log file: %v\n", time.Now().Format("2006/01/02 15:04:05"))))
		}
	}

	if _, err = hook.file.Write(output); err != nil {
		return fmt.Errorf("Unable to print logs to file: %w", err)
	}

	return nil
}

var fileMutex sync.Mutex
