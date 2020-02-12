package multilogger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/sirupsen/logrus"
)

func cleanupModuleName(moduleName string) string {
	return strings.Trim(strings.Map(
		func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '/' || r == ':' {
				return r
			}
			return -1
		},
		moduleName,
	), "/")
}

type fileHook struct {
	*genericHook
	filename  string
	isDir     bool
	file      *os.File
	addHeader bool
}

func (hook *fileHook) clone() logrus.Hook {
	// Duplicate the file hook to ensure that the copy
	// has its own attributes when the object is copied.
	return &fileHook{
		genericHook: hook.genericHook.clone(),
		filename:    hook.filename,
		isDir:       hook.isDir,
		file:        hook.file,
		addHeader:   hook.addHeader,
	}
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
		targetFile := hook.filename
		if hook.isDir {
			moduleName := cleanupModuleName(hook.logger.GetModule())
			targetFile = path.Join(hook.filename, strings.Replace(moduleName, ":", ".", -1)) + ".log"
			if targetFile, err = filepath.Abs(targetFile); err != nil {
				return err
			}
			if hook.file != nil && hook.file.Name() != targetFile {
				hook.file = nil
			}
		}
		if hook.file == nil {
			logDir := path.Dir(targetFile)
			logFileExists := false
			if _, err := os.Stat(logDir); os.IsNotExist(err) {
				// Log directory doesn't exist, create it
				if err := os.MkdirAll(logDir, 0777); err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
				if err := os.Chmod(logDir, 0777); err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
			} else {
				if _, err := os.Stat(targetFile); err == nil {
					logFileExists = true
				}
			}

			if hook.file, err = os.OpenFile(targetFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			if err := os.Chmod(targetFile, 0777); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			if hook.addHeader {
				if logFileExists {
					// Add a bit of whitespace before logging
					if err := hook.printf(name, hook.file, "\n"); err != nil {
						return fmt.Errorf("%s: %w", name, err)
					}
				}
				if err := hook.printf(name, hook.file, "# %v\n", entry.Time.Format(defaultTimestampFormat)); err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
			}
		}

		return hook.printf(name, hook.file, string(output))
	})
}

var fileMutex sync.Mutex
