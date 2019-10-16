package multilogger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLoggingLevels(t *testing.T) {
	tests := []struct {
		name         string
		minimumLevel logrus.Level
		levels       []logrus.Level
	}{
		{
			name:         "disabled",
			minimumLevel: DisabledLevel,
			levels:       nil,
		},
		{
			name:         "panic",
			minimumLevel: logrus.PanicLevel,
			levels:       []logrus.Level{logrus.PanicLevel},
		},
		{
			name:         "fatal",
			minimumLevel: logrus.FatalLevel,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel},
		},
		{
			name:         "error",
			minimumLevel: logrus.ErrorLevel,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
		},
		{
			name:         "warn",
			minimumLevel: logrus.WarnLevel,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel},
		},
		{
			name:         "info",
			minimumLevel: logrus.InfoLevel,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel},
		},
		{
			name:         "debug",
			minimumLevel: logrus.DebugLevel,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel},
		},
		{
			name:         "trace",
			minimumLevel: logrus.TraceLevel,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel},
		},
		{
			name:         "over",
			minimumLevel: 7,
			levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel, logrus.TraceLevel, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &GenericHook{MinimumLevel: tt.minimumLevel}
			assert.Equal(t, tt.levels, hook.Levels())
		})
	}
}
