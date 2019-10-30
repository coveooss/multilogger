package multilogger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// NewHook generates a named hook wrapper that is able the handle its own logging level.
//
// level: Accept any kind of object, but must be resolvable into a valid logrus level name.
func NewHook(name string, level interface{}, hook logrus.Hook) *Hook {
	return &Hook{
		name:  name,
		level: ParseLogLevel(level),
		inner: hook,
	}
}

// NewConsoleHook creates a new hook to log information to console (default to stderr).
//
// level: Accept any kind of object, but must be resolvable into a valid logrus level name.
func NewConsoleHook(name string, level interface{}, format ...interface{}) *Hook {
	if name == "" {
		name = consoleHookName
	}
	return NewHook(name, level, &consoleHook{
		out:         os.Stdout,
		log:         logrus.StandardLogger().Out,
		genericHook: &genericHook{formatter: getFormatter(true, format...)},
	})
}

// NewFileHook creates a new hook to log information into a file.
//
// level: Accept any kind of object, but must be resolvable into a valid logrus level name.
func NewFileHook(filename string, level interface{}, format ...interface{}) *Hook {
	if format == nil {
		format = append(format, NewFormatter(false, os.Getenv(FormatFileEnvVar), os.Getenv(FormatEnvVar), DefaultFileFormat))
	}
	return NewHook(filename, level, &fileHook{
		genericHook: &genericHook{formatter: getFormatter(false, format...)},
		filename:    filename,
		addHeader:   true,
	})
}

func getFormatter(color bool, format ...interface{}) (formatter logrus.Formatter) {
	switch len(format) {
	case 0:
		return NewFormatter(color, os.Getenv(FormatEnvVar), DefaultConsoleFormat)
	case 1:
		if f, ok := format[0].(logrus.Formatter); ok {
			return f
		}
	}
	return NewFormatter(color, format...)
}

// Hook represents a hook that logs at a given level. This struct must be extended to implement the Fire func.
type Hook struct {
	name  string
	inner logrus.Hook
	level logrus.Level
}

// Levels returns the levels that should be handled by the hook.
func (hook *Hook) Levels() []logrus.Level {
	if hook.inner == nil {
		return nil
	}

	if levels := hook.inner.Levels(); levels == nil {
		result := make([]logrus.Level, 0, hook.level+1)
		if isGeneric, _ := hook.inner.(genericHookI); isGeneric != nil {
			result = append(result, outputLevel)
		}
		// We compute the array of levels based on specified level
		switch level := hook.level; {
		case level == DisabledLevel:
			return result
		case level <= logrus.TraceLevel:
			return append(result, logrus.AllLevels[:level+1]...)
		default:
			result := append(result, logrus.AllLevels...)
			for i := logrus.TraceLevel; i < level; i++ {
				result = append(result, i+1)
			}
			return result
		}
	} else {
		// We limit the actual hook levels to the maximum accepted level
		result := make([]logrus.Level, 0, len(levels))
		for _, level := range levels {
			if level <= hook.level {
				result = append(result, level)
			}
		}
		return result
	}
}

// Fire writes logs to the console.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	if hook.inner == nil {
		return fmt.Errorf("Hook not configured properly")
	}
	return hook.inner.Fire(entry)
}

// GetFormatter returns the formater associated to the hook.
// The function will panic if called upon a hook that do not support formatter.
func (hook *Hook) GetFormatter() logrus.Formatter {
	return hook.inner.(genericHookI).Formatter()
}

// Formatter returns the Formatter associated to the hook.
// The function will panic if called upon a hook that do not support formatter.
func (hook *Hook) Formatter() *Formatter {
	f, _ := hook.GetFormatter().(*Formatter)
	return f
}

// SetFormatter allows setting a formatter on hook that support it.
// The function will panic if called upon a hook that do not support formatter.
func (hook *Hook) SetFormatter(formatter logrus.Formatter) *Hook {
	hook.inner.(genericHookI).SetFormatter(formatter)
	return hook
}

// SetFormat allows setting a format string on hook that support it.
// The function will panic if called upon a hook that do not support formatter.
func (hook *Hook) SetFormat(formats ...interface{}) *Hook {
	hook.Formatter().SetLogFormat(formats...)
	return hook
}

// SetColor allows setting color mode on hook that support it.
// The function will panic if called upon a hook that do not support formatter.
func (hook *Hook) SetColor(color bool) *Hook {
	if f := hook.Formatter(); f != nil {
		f.SetColor(color)
	}
	if f, _ := hook.GetFormatter().(*logrus.TextFormatter); f != nil {
		f.ForceColors = color
		f.DisableColors = !color
	}
	return hook
}

// SetOut allows configuring the logging stream for a hook that support it.
// The function will panic if called upon a hook that is not a valid ConsoleHook.
func (hook *Hook) SetOut(out io.Writer) *Hook {
	hook.inner.(consoleI).SetOut(out)
	return hook
}

// SetStdout allows configuring the output stream for a hook that support it.
// The function will panic if called upon a hook that is not a valid ConsoleHook.
func (hook *Hook) SetStdout(out io.Writer) *Hook {
	hook.inner.(consoleI).SetStdout(out)
	return hook
}
