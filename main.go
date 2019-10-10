package multilogger

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	timestampFormat = "2006/01/02 15:04:05.000"
	// DisabledLevel can be set when one of the logging hooks should be disabled
	DisabledLevel logrus.Level = 9999
	// DisabledLevelName is the textual name of the
	DisabledLevelName string = "disabled"
)

// AcceptedLevels returns all accepted logrus levels
func AcceptedLevels() []string {
	levels := []string{DisabledLevelName}
	for _, level := range logrus.AllLevels {
		levels = append(levels, level.String())
	}
	return levels
}

// AcceptedLevelsString returns all accepted logrus levels as a comma-separated string
func AcceptedLevelsString() string {
	return strings.Join(AcceptedLevels(), ", ")
}

// MustParseLogLevel converts a string or number into a logging level
// It panics if the supplied valid cannot be converted into a valid logrus Level
func MustParseLogLevel(level interface{}) logrus.Level {
	result, err := ParseLogLevel(level)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseLogLevel converts a string or number into a logging level
func ParseLogLevel(level interface{}) (logrus.Level, error) {
	logLevelFromInt := func(levelNum int) logrus.Level {
		if levelNum == int(DisabledLevel) {
			return DisabledLevel
		}
		return logrus.Level(levelNum)
	}

	if logrusLevel, ok := level.(logrus.Level); ok {
		return logrusLevel, nil
	}

	if levelNum, ok := level.(int); ok {
		return logLevelFromInt(levelNum), nil
	}

	if levelString, ok := level.(string); ok {
		if levelNum, err := strconv.Atoi(levelString); err == nil {
			return logLevelFromInt(levelNum), nil
		}

		if strings.ToLower(levelString) == DisabledLevelName {
			return DisabledLevel, nil
		}
		parsedLevel, err := logrus.ParseLevel(levelString)
		if err != nil {
			return DisabledLevel, fmt.Errorf("Unable to parse logging level: %v", err)
		}
		return parsedLevel, nil
	}
	return DisabledLevel, fmt.Errorf("Unable to parse the given logging level %v. It has to be a string or an integer", level)
}

// MultiLogger represents a logger that logs to both a file and the console at different (configurable) levels
type MultiLogger struct {
	*logrus.Logger

	// Loggers
	consoleHook *ConsoleHook
	fileHook    *FileHook
}

// New creates a new Multilogger instance and panic if it can't be created
func New(consoleLevel interface{}, fileLevel interface{}, filename string, module string) *MultiLogger {
	logger, err := Create(consoleLevel, fileLevel, filename, module)
	if err != nil {
		panic(err)
	}
	return logger
}

// Create creates a new Multilogger instance
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

// SetConsoleLevel modifies the logging level for the console logger
func (logger *MultiLogger) SetConsoleLevel(level interface{}) {
	logger.consoleHook.MinimumLevel = MustParseLogLevel(level)
	logger.refreshLoggers()
}

// SetFileLevel modifies the logging level for the file logger
func (logger *MultiLogger) SetFileLevel(level interface{}) {
	logger.fileHook.MinimumLevel = MustParseLogLevel(level)
	logger.refreshLoggers()
}

// ConfigureFileLogger modifies the logging level and filename for the file logger
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
