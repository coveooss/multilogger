package multilogger

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	// DisabledLevel can be set when one of the logging hooks should be disabled
	DisabledLevel     logrus.Level = math.MaxUint32
	disabledLevelName string       = "disabled"
)

// AcceptedLevels returns all accepted logrus levels.
func AcceptedLevels() []string {
	levels := []string{disabledLevelName}
	for _, level := range logrus.AllLevels {
		levels = append(levels, level.String())
	}
	return levels
}

// AcceptedLevelsString returns all accepted logrus levels as a comma-separated string.
func AcceptedLevelsString() string {
	return strings.Join(AcceptedLevels(), ", ")
}

// MustParseLogLevel converts a string or number into a logging level.
// It panics if the supplied valid cannot be converted into a valid logrus Level.
func MustParseLogLevel(level interface{}) logrus.Level {
	result, err := ParseLogLevel(level)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseLogLevel converts a string or number into a logging level.
func ParseLogLevel(level interface{}) (logrus.Level, error) {
	logLevelFromInt := func(levelNum int) logrus.Level {
		if levelNum == int(DisabledLevel) {
			return DisabledLevel
		}
		return logrus.Level(levelNum)
	}

	if level == nil || level == "" {
		return DisabledLevel, nil
	} else if logrusLevel, ok := level.(logrus.Level); ok {
		return logrusLevel, nil
	} else if levelNum, ok := level.(int); ok {
		return logLevelFromInt(levelNum), nil
	} else if levelString, ok := level.(string); ok {
		if levelNum, err := strconv.Atoi(levelString); err == nil {
			return logLevelFromInt(levelNum), nil
		}

		if strings.ToLower(levelString) == disabledLevelName {
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
