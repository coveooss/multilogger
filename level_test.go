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
	assert.Equal(t, DisabledLevel, ParseLogLevel(nil))
	assert.Equal(t, DisabledLevel, ParseLogLevel(""))
	assert.Equal(t, DisabledLevel, ParseLogLevel(DisabledLevel))
	assert.Equal(t, DisabledLevel, ParseLogLevel(disabledLevelName))
	assert.Equal(t, DisabledLevel, ParseLogLevel(int(DisabledLevel)))
	assert.Equal(t, DisabledLevel, ParseLogLevel(strconv.Itoa(int(DisabledLevel))))
	for _, level := range logrus.AllLevels {
		name := level.String()
		assert.Equal(t, level, ParseLogLevel(level))
		assert.Equal(t, level, ParseLogLevel(name))
		assert.Equal(t, level, ParseLogLevel(strings.ToUpper(name)))
		assert.Equal(t, level, ParseLogLevel(strings.ToLower(name)))
		assert.Equal(t, level, ParseLogLevel(strings.ToTitle(name)))
		assert.Equal(t, level, ParseLogLevel(int(level)))
		type derivedType string
		assert.Equal(t, level, ParseLogLevel(derivedType(fmt.Sprint(int(level)))))
		assert.Equal(t, level, ParseLogLevel(derivedType(name)))
	}
	assert.Panics(t, func() { ParseLogLevel("Invalid") })
	assert.Panics(t, func() { ParseLogLevel(1.23) })
}

func TestGetAcceptedLevels(t *testing.T) {
	assert.Equal(t, "disabled, panic, fatal, error, warning, info, debug, trace", AcceptedLevelsString())
}

func TestParseInvalidLogLevel(t *testing.T) {
	level, err := TryParseLogLevel("invalid")
	assert.Equal(t, level, DisabledLevel)
	assert.EqualError(t, err, `Unable to parse logging level: not a valid logrus Level: "invalid"`)
	level, err = TryParseLogLevel(1.234)
	assert.Equal(t, level, DisabledLevel)
	assert.EqualError(t, err, `Unable to parse logging level: not a valid logrus Level: "1.234"`)
}
