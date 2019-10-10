package multilogger

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// GenericHook represents a hook that logs at a given level. This struct must be extended to implement the Fire func
type GenericHook struct {
	Formatter    logrus.Formatter
	MinimumLevel logrus.Level
}

func (hook *GenericHook) formatEntry(entry *logrus.Entry) ([]byte, error) {
	if hook.MinimumLevel == DisabledLevel {
		return []byte{}, nil
	}
	formatted, err := hook.Formatter.Format(entry)
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to format the given log entry: %v", err)
	}
	return formatted, nil
}

// Levels returns the levels that should be handled by the hook
func (hook *GenericHook) Levels() []logrus.Level {
	levels := []logrus.Level{}
	for _, level := range logrus.AllLevels {
		levels = append(levels, level)
		if level == hook.MinimumLevel {
			return levels
		}
	}
	return []logrus.Level{}
}

// GetColor returns an ANSI color formatting function for every logrus logging level
func GetColor(level logrus.Level) func(format string, args ...interface{}) string {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return color.HiGreenString
	case logrus.InfoLevel:
		return color.HiBlueString
	case logrus.WarnLevel:
		return color.YellowString
	case logrus.ErrorLevel, logrus.FatalLevel:
		return color.RedString
	case logrus.PanicLevel:
		return color.MagentaString
	}
	return fmt.Sprintf
}
