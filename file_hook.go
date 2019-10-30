package multilogger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type fileHook struct {
	*genericHook
	filename  string
	file      *os.File
	addHeader bool
}

func (hook *fileHook) Fire(entry *logrus.Entry) (err error) {
	return hook.fire(entry, func() error {
		name := fmt.Sprintf("FileHook %s", hook.filename)
		output := entry.Message
		if entry.Level != outputLevel {
			if output, err = hook.formatEntry(name, entry); err != nil {
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
				return fmt.Errorf("%s: %w", name, err)
			}
			if hook.addHeader {
				if logFileExists {
					// Add a bit of whitespace before logging
					if err := hook.printf(name, hook.file, "\n"); err != nil {
						return err
					}
				}
				if err := hook.printf(name, hook.file, "# %v\n", entry.Time.Format(defaultTimestampFormat)); err != nil {
					return err
				}
			}
		}

		return hook.printf(name, hook.file, string(output))
	})
}

var fileMutex sync.Mutex
