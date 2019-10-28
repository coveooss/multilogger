package multilogger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLoggingLevels(t *testing.T) {
	tests := []struct {
		name   string
		level  logrus.Level
		levels []logrus.Level
	}{
		{"disabled", DisabledLevel, []logrus.Level{outputLevel}},
		{"panic", logrus.PanicLevel, []logrus.Level{outputLevel, logrus.PanicLevel}},
		{"fatal", logrus.FatalLevel, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel}},
		{"error", logrus.ErrorLevel, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}},
		{"warn", logrus.WarnLevel, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}},
		{"info", logrus.InfoLevel, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel}},
		{"debug", logrus.DebugLevel, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel}},
		{"trace", logrus.TraceLevel, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel}},
		{"over", 7, []logrus.Level{outputLevel, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := NewHook("", tt.level, &genericHook{})
			assert.Equal(t, tt.levels, hook.Levels())
		})
	}
}
