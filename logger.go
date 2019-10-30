package multilogger

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// CallerEnvVar is an environment variable that enable the caller stack by default.
	CallerEnvVar = "MULTILOGGER_CALLER"
	// FormatEnvVar is an environment variable that allows users to set the default format used for log entry.
	FormatEnvVar = "MULTILOGGER_FORMAT"
	// FormatFileEnvVar is an environment variable that allows users to set the default format used for log entry using a file logger.
	FormatFileEnvVar = "MULTILOGGER_FILE_FORMAT"
	// DefaultConsoleFormat is the format used by NewConsoleHook if MULTILOGGER_FORMAT is not set.
	DefaultConsoleFormat = "%module:Italic,Green,SquareBrackets,IgnoreEmpty,Space%%time% %-8level:upper,color% %message:color%"
	// DefaultFileFormat is the format used by NewFileHook if neither MULTILOGGER_FORMAT or MULTILOGGER_FILE_FORMAT are set.
	DefaultFileFormat = "[%module,SquareBrackets,IgnoreEmpty,Space%]%time% %-8level:upper% %message%"
)

const (
	consoleHookName = "console-hook"
	moduleFieldName = "module-field"
)

// Logger represents a logger that logs to both a file and the console at different (configurable) levels.
type Logger struct {
	*logrus.Entry
	hooks      map[string]*leveledHook
	level      logrus.Level
	PrintLevel logrus.Level
}

type leveledHook struct {
	level logrus.Level
	hook  *Hook
}

// New creates a new Multilogger instance.
// If no hook is provided, it defaults to standard console logger at warning log level.
func New(module string, hooks ...*Hook) *Logger {
	if len(hooks) == 0 {
		hooks = []*Hook{NewConsoleHook("", logrus.WarnLevel)}
	}
	logger := &Logger{Entry: logrus.New().WithField(moduleFieldName, module)}
	if caller := os.Getenv(CallerEnvVar); caller != "" {
		caller, err := strconv.ParseBool(caller)
		logger.SetReportCaller(err != nil || caller)
	}
	logger.Logger.Out = ioutil.Discard      // Discard all logs to the main logger
	logger.Logger.Level = DisabledLevel - 1 // Always log at the highest possible level. Hooks will decide if the log goes through or not
	logger.AddHooks(hooks...)
	logger.PrintLevel = outputLevel
	return logger
}

// Copy returns a new logger with the same hooks but a different module name.
func (logger *Logger) Copy(module string) *Logger {
	var hooks []*Hook
	for key, hook := range logger.hooks {
		hooks = append(hooks, NewHook(key, hook.level, hook.hook.inner))
	}
	newLogger := New(module, hooks...)
	newLogger.Entry = logger.Entry.WithField(moduleFieldName, module)
	newLogger.PrintLevel = logger.PrintLevel
	return newLogger
}

// WithTime return a new logger with a fixed time for log entry (useful for testing).
func (logger *Logger) WithTime(time time.Time) *Logger {
	newLogger := logger.Copy(logger.GetModule())
	newLogger.Entry = newLogger.Entry.WithTime(time)
	return newLogger
}

// AddTime add the specified duration to the current logger if its time has been freezed. Useful for testing.
func (logger *Logger) AddTime(duration time.Duration) *Logger {
	if !logger.Entry.Time.IsZero() {
		logger.Entry.Time = logger.Entry.Time.Add(duration)
	}
	return logger
}

// WithField return a new logger with a single additional entry.
func (logger *Logger) WithField(key string, value interface{}) *Logger {
	return logger.WithFields(logrus.Fields{key: value})
}

// WithFields return a new logger with a new fields value.
func (logger *Logger) WithFields(fields logrus.Fields) *Logger {
	newLogger := logger.Copy(logger.GetModule())
	newLogger.Entry = logger.Entry.WithFields(fields)
	return newLogger
}

// WithContext return a new logger with a new context.
func (logger *Logger) WithContext(ctx context.Context) *Logger {
	newLogger := logger.Copy(logger.GetModule())
	newLogger.Entry = logger.Entry.WithContext(ctx)
	return newLogger
}

// Print acts as fmt.Print but sends the output to a special logging level that allows multiple output support through Hooks.
//
// ATTENTION, default behaviour for logrus.Print is to log at Info level.
func (logger *Logger) Print(args ...interface{}) {
	logger.Entry.Log(logger.PrintLevel, args...)
}

// Println acts as fmt.Println but sends the output to a special logging level that allows multiple output support through Hooks.
//
// ATTENTION, default behaviour for logrus.Println is to log at Info level.
func (logger *Logger) Println(args ...interface{}) {
	logger.Print(fmt.Sprintln(args...))
}

// Printf acts as fmt.Printf but sends the output to a special logging level that allows multiple output support through Hooks.
//
// ATTENTION, default behaviour for logrus.Printf is to log at Info level.
func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.Entry.Logf(logger.PrintLevel, format, args...)
}

// GetLevel returns the highest logger level registered by the hooks.
func (logger *Logger) GetLevel() logrus.Level {
	if mainLevel := logger.Logger.GetLevel(); mainLevel < logger.level {
		return mainLevel
	}
	return logger.level
}

// GetModule returns the module name associated to the current logger.
func (logger *Logger) GetModule() string {
	return logger.Data[moduleFieldName].(string)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level logrus.Level) bool {
	return logger.GetLevel() >= level
}

// SetReportCaller enables caller reporting to be added to each log entry.
func (logger *Logger) SetReportCaller(reportCaller bool) *Logger {
	logger.Logger.SetReportCaller(reportCaller)
	return logger
}

// SetExitFunc let user define what should be executed when a logging call exit (default is call to os.Exit(int)).
func (logger *Logger) SetExitFunc(exitFunc func(int)) {
	logger.Logger.ExitFunc = exitFunc
}

// AddHook adds a hook to the hook collection and associated it with a name and a level.
// Can also be used to replace an existing hook.
func (logger *Logger) AddHook(name string, level interface{}, hook logrus.Hook) *Logger {
	return logger.AddHooks(NewHook(name, level, hook))
}

// AddHooks adds a collection of hook wrapper as hook to the current logger.
func (logger *Logger) AddHooks(hooks ...*Hook) *Logger {
	if logger.hooks == nil {
		logger.hooks = make(map[string]*leveledHook)
	}
	for _, hook := range hooks {
		logger.hooks[hook.name] = &leveledHook{hook.level, hook}
	}
	return logger.refreshLoggers()
}

// AddConsole adds a console hook to the current logger.
func (logger *Logger) AddConsole(name string, level interface{}, format ...interface{}) *Logger {
	return logger.AddHooks(NewConsoleHook(name, level, format...))
}

// AddFile adds a file hook to the current logger.
func (logger *Logger) AddFile(filename string, level interface{}, format ...interface{}) *Logger {
	return logger.AddHooks(NewFileHook(filename, level, format...))
}

// RemoveHook deletes a hook from the hook collection.
func (logger *Logger) RemoveHook(name string) *Logger {
	delete(logger.hooks, name)
	return logger.refreshLoggers()
}

// Hook returns the hook identified by name.
func (logger *Logger) Hook(name string) *Hook {
	if hook := logger.getHook(name); hook != nil {
		return hook.hook
	}
	return nil
}

// GetHookLevel returns the logging level associated with a specific logger.
func (logger *Logger) GetHookLevel(name string, level interface{}) logrus.Level {
	if hook := logger.getHook(name); hook != nil {
		return hook.level
	}
	return DisabledLevel
}

// SetHookLevel set a new log level for a registered hook.
func (logger *Logger) SetHookLevel(name string, level interface{}) error {
	if hook := logger.Hook(name); hook != nil {
		logger.AddHook(hook.name, level, hook.inner)
	}
	return fmt.Errorf("Hook not found %s", name)
}

// ListHooks returns the list of registered hook names.
func (logger *Logger) ListHooks() []string {
	result := make([]string, 0, len(logger.hooks))
	for key := range logger.hooks {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func (logger *Logger) refreshLoggers() *Logger {
	logger.Logger.Hooks = make(logrus.LevelHooks)
	var level logrus.Level
	for _, hook := range logger.hooks {
		logger.Logger.Hooks.Add(hook.hook)
		if hook.level > level {
			level = hook.level
		}
	}
	logger.level = level
	return logger
}

func (logger *Logger) getHook(name string) *leveledHook {
	if name == "" {
		name = consoleHookName
	}
	return logger.hooks[name]
}

// Writer is the implementation of io.Writer. You should not call directly this function.
// The function will fail if called directly on a stream that have not been configured as out stream.
func (logger *Logger) Write(p []byte) (n int, err error) {
	logger.Print(string(p))
	return len(p), nil
}
