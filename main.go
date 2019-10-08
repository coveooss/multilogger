package multilogger

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

// DisabledLevel can be set when one of the logging hooks should be disabled
const DisabledLevel = 9999

// MultiLogger represents a logger that logs to both a file and the console at different (configurable) levels
type MultiLogger struct {
	*logrus.Logger

	// Loggers
	consoleHook *ConsoleHook
	fileHook    *FileHook
}

// New creates a new Multilogger instance
func New(consoleLevel logrus.Level, fileLevel logrus.Level, fileName string, module string) *MultiLogger {
	logger := &MultiLogger{
		Logger: logrus.New(),
	}
	logger.Out = ioutil.Discard      // Discard all logs to the main logger
	logger.Level = logrus.TraceLevel // Always log at TRACE level. Hooks will decide if the log goes through or not

	formatter := &easy.Formatter{
		TimestampFormat: timestampFormat,
		LogFormat:       fmt.Sprintf("[%s]", module) + " %time% %lvl% %msg%\n",
	}
	if fileName == "" {
		fileName = module + ".log"
	}
	logger.fileHook = NewFileHook(fileName, fileLevel, formatter)
	logger.consoleHook = NewConsoleHook(consoleLevel, formatter)
	logger.refreshLoggers()

	return logger
}

// SetConsoleLevel modifies the logging level for the console logger
func (logger *MultiLogger) SetConsoleLevel(level logrus.Level) {
	logger.consoleHook.MinimumLevel = level
	logger.refreshLoggers()
}

// SetFileLevel modifies the logging level for the file logger
func (logger *MultiLogger) SetFileLevel(level logrus.Level) {
	logger.fileHook.MinimumLevel = level
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
