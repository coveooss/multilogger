package multilogger

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/sirupsen/logrus"
)

const timestampFormat = "2006/01/02 15:04:05.000"

// MultiLogger represents a logger that logs to both a file and the console at different (configurable) levels.
type MultiLogger struct {
	*logrus.Logger

	// Loggers
	consoleHook *ConsoleHook
	fileHook    *FileHook
}

// New creates a new Multilogger instance and panic if it can't be created.
func New(consoleLevel interface{}, fileLevel interface{}, filename string, module string) *MultiLogger {
	logger, err := Create(consoleLevel, fileLevel, filename, module)
	if err != nil {
		panic(err)
	}
	return logger
}

// Create creates a new Multilogger instance.
func Create(consoleLevel interface{}, fileLevel interface{}, filename string, module string) (logger *MultiLogger, err error) {
	if consoleLevel, err = ParseLogLevel(consoleLevel); err != nil {
		return
	}
	if fileLevel, err = ParseLogLevel(fileLevel); err != nil {
		return
	}
	logger = &MultiLogger{Logger: logrus.New()}
	logger.Out = ioutil.Discard      // Discard all logs to the main logger
	logger.Level = logrus.TraceLevel // Always log at TRACE level. Hooks will decide if the log goes through or not

	format := fmt.Sprintf("[%s]", module) + " %time% %lvl% %msg%\n"
	logger.fileHook = NewFileHook(filename, fileLevel.(logrus.Level), &Formatter{
		TimestampFormat: timestampFormat,
		LogFormat:       format,
		Color:           false,
	})
	logger.consoleHook = NewConsoleHook(consoleLevel.(logrus.Level), &Formatter{
		TimestampFormat: timestampFormat,
		LogFormat:       format,
		Color:           true,
	})
	logger.refreshLoggers()

	return
}

// GetLevel returns the logger level.
func (logger *MultiLogger) GetLevel() logrus.Level {
	return logrus.Level(math.Max(float64(logger.consoleHook.MinimumLevel), float64(logger.fileHook.MinimumLevel)))
}

// SetConsoleLevel modifies the logging level for the console logger.
func (logger *MultiLogger) SetConsoleLevel(level interface{}) {
	logger.consoleHook.MinimumLevel = MustParseLogLevel(level)
	logger.refreshLoggers()
}

// SetFileLevel modifies the logging level for the file logger.
func (logger *MultiLogger) SetFileLevel(level interface{}) {
	logger.fileHook.MinimumLevel = MustParseLogLevel(level)
	logger.refreshLoggers()
}

// ConfigureFileLogger modifies the logging level and filename for the file logger.
func (logger *MultiLogger) ConfigureFileLogger(level interface{}, newFilename string) {
	logger.fileHook.MinimumLevel = MustParseLogLevel(level)
	logger.fileHook.SetFilename(newFilename)
	logger.refreshLoggers()
}

func (logger *MultiLogger) refreshLoggers() {
	logger.Hooks = make(logrus.LevelHooks)
	if logger.consoleHook != nil {
		logger.Hooks.Add(logger.consoleHook)
	}
	if logger.fileHook != nil {
		logger.Hooks.Add(logger.fileHook)
	}
}
