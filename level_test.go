package multilogger

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	assert.Equal(t, DisabledLevel, MustParseLogLevel(nil))
	assert.Equal(t, DisabledLevel, MustParseLogLevel(""))
	assert.Equal(t, DisabledLevel, MustParseLogLevel(DisabledLevel))
	assert.Equal(t, DisabledLevel, MustParseLogLevel(disabledLevelName))
	assert.Equal(t, DisabledLevel, MustParseLogLevel(int(DisabledLevel)))
	assert.Equal(t, DisabledLevel, MustParseLogLevel(strconv.Itoa(int(DisabledLevel))))
	for _, level := range logrus.AllLevels {
		name := level.String()
		assert.Equal(t, level, MustParseLogLevel(level))
		assert.Equal(t, level, MustParseLogLevel(name))
		assert.Equal(t, level, MustParseLogLevel(strings.ToUpper(name)))
		assert.Equal(t, level, MustParseLogLevel(strings.ToLower(name)))
		assert.Equal(t, level, MustParseLogLevel(strings.ToTitle(name)))
		assert.Equal(t, level, MustParseLogLevel(int(level)))
		assert.Equal(t, level, MustParseLogLevel(strconv.Itoa(int(level))))
	}
	assert.Panics(t, func() { MustParseLogLevel("Invalid") })
	assert.Panics(t, func() { MustParseLogLevel(1.23) })
}

func TestGetAcceptedLevels(t *testing.T) {
	assert.Equal(t, "disabled, panic, fatal, error, warning, info, debug, trace", AcceptedLevelsString())
}

func TestParseInvalidLogLevel(t *testing.T) {
	level, err := ParseLogLevel("invalid")
	assert.Equal(t, level, DisabledLevel)
	assert.EqualError(t, err, `Unable to parse logging level: not a valid logrus Level: "invalid"`)
	level, err = ParseLogLevel(1.234)
	assert.Equal(t, level, DisabledLevel)
	assert.EqualError(t, err, `Unable to parse the given logging level 1.234. It has to be a string or an integer`)
}

func xExample() {
	logger, _ := Create(7, "debug", "", "test")
	logger.SetReportCaller(true)
	logger.Warn("Hello1")
	logger.ConfigureFileLogger("debug", "tata")
	logger.Error("This is an error message")
	logger.Warn("This is a warning message")
	logger.Info("This is an information message")
	logger.Debug("This is a debug message")
	logger.Trace("This is a trace message")
	logger.WithField("hello", 1).Warn("This message has context")
	fmt.Println(logger.consoleHook.Levels())
	// Output:
}
