package multilogger

import (
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	assert.Equal(t, DisabledLevel, ParseLogLevel(DisabledLevel))
	assert.Equal(t, DisabledLevel, ParseLogLevel(DisabledLevelName))
	assert.Equal(t, DisabledLevel, ParseLogLevel(int(DisabledLevel)))
	assert.Equal(t, DisabledLevel, ParseLogLevel(strconv.Itoa(int(DisabledLevel))))
	for _, level := range logrus.AllLevels {
		assert.Equal(t, level, ParseLogLevel(level))
		assert.Equal(t, level, ParseLogLevel(level.String()))
		assert.Equal(t, level, ParseLogLevel(int(level)))
		assert.Equal(t, level, ParseLogLevel(strconv.Itoa(int(level))))
	}
}

func TestGetAcceptedLevels(t *testing.T) {
	assert.Equal(t, "disabled, panic, fatal, error, warning, info, debug, trace", AcceptedLevelsString())
}
