package multilogger

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

// DisabledLevel can be set when one of the logging hooks should be disabled
const DisabledLevel logrus.Level = 9999

// DisabledLevelName is the textual name of the
const DisabledLevelName string = "disabled"

func AcceptedLevels() []string {
	levels := []string{DisabledLevelName}
	for _, level := range logrus.AllLevels {
		levels = append(levels, level.String())
	}
	return levels
}

func AcceptedLevelsString() string {
	return strings.Join(AcceptedLevels(), ", ")
}

// ParseLogLevel converts a string or number into a logging level
func ParseLogLevel(level interface{}) logrus.Level {
	logLevelFromInt := func(levelNum int) logrus.Level {
		if levelNum == int(DisabledLevel) {
			return DisabledLevel
		}
		return logrus.Level(levelNum)
	}

	if logrusLevel, ok := level.(logrus.Level); ok {
		return logrusLevel
	}

	if levelNum, ok := level.(int); ok {
		return logLevelFromInt(levelNum)
	}

	if levelString, ok := level.(string); ok {
		if levelNum, err := strconv.Atoi(levelString); err == nil {
			return logLevelFromInt(levelNum)
		}

		if strings.ToLower(levelString) == DisabledLevelName {
			return DisabledLevel
		}
		parsedLevel, err := logrus.ParseLevel(levelString)
		if err != nil {
			panic(fmt.Errorf("Unable to parse the given logging level %s: %v", level, err))
		}
		return parsedLevel
	}
	panic(fmt.Errorf("Unable to parse the given logging level %v. It has to be a string or an integer", level))
}

// MultiLogger represents a logger that logs to both a file and the console at different (configurable) levels
type MultiLogger struct {
	*logrus.Logger

	// Loggers
	consoleHook *ConsoleHook
	fileHook    *FileHook
}

// New creates a new Multilogger instance
func New(consoleLevel interface{}, fileLevel interface{}, fileName string, module string) *MultiLogger {
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
	logger.fileHook = NewFileHook(fileName, ParseLogLevel(fileLevel), formatter)
	logger.consoleHook = NewConsoleHook(ParseLogLevel(consoleLevel), formatter)
	logger.refreshLoggers()

	return logger
}

// SetConsoleLevel modifies the logging level for the console logger
func (logger *MultiLogger) SetConsoleLevel(level interface{}) {
	logger.consoleHook.MinimumLevel = ParseLogLevel(level)
	logger.refreshLoggers()
}

// SetFileLevel modifies the logging level for the file logger
func (logger *MultiLogger) SetFileLevel(level interface{}) {
	logger.fileHook.MinimumLevel = ParseLogLevel(level)
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