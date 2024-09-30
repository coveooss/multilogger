package multilogger

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/coveooss/multilogger/errors"
	"github.com/sirupsen/logrus"
)

const (
	// DisabledLevel can be set when one of the logging hooks should be disabled.
	DisabledLevel     logrus.Level = math.MaxUint32
	outputLevel                    = DisabledLevel - 1
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

// ParseLogLevel converts a string or number into a logging level.
// It panics if the supplied valid cannot be converted into a valid logrus Level.
func ParseLogLevel(level interface{}) logrus.Level {
	return errors.Must(TryParseLogLevel(level)).(logrus.Level)
}

// TryParseLogLevel converts a string or number into a logging level.
func TryParseLogLevel(level interface{}) (logrus.Level, error) {
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
	}
	levelString := fmt.Sprint(level)
	if levelNum, err := strconv.Atoi(levelString); err == nil {
		return logLevelFromInt(levelNum), nil
	}
	if strings.ToLower(levelString) == disabledLevelName {
		return DisabledLevel, nil
	}
	parsedLevel, err := logrus.ParseLevel(levelString)
	if err != nil {
		return DisabledLevel, fmt.Errorf("unable to parse logging level: %w", err)
	}
	return parsedLevel, nil
}
