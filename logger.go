package multilogger

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/coveooss/multilogger/errors"
	"github.com/coveooss/multilogger/reutils"
	"github.com/sirupsen/logrus"
)

const (
	// CallerEnvVar is an environment variable that enable the caller stack by default.
	CallerEnvVar = "MULTILOGGER_CALLER"
	// FormatEnvVar is an environment variable that allows users to set the default format used for log entry.
	FormatEnvVar = "MULTILOGGER_FORMAT"
	// FormatFileEnvVar is an environment variable that allows users to set the default format used for log entry using a file logger.
	FormatFileEnvVar = "MULTILOGGER_FILE_FORMAT"
	// DefaultFileFormat is the format used by NewFileHook if neither MULTILOGGER_FORMAT or MULTILOGGER_FILE_FORMAT are set.
	DefaultFileFormat = "%module:SquareBrackets,IgnoreEmpty,Space%%time% %-8level:upper% %message%"
)

const (
	moduleFieldName = "module-field"
)

// Logger represents a logger that logs to both a file and the console at different (configurable) levels.
type Logger struct {
	*logrus.Entry
	PrintLevel logrus.Level
	Catcher    bool

	hooks     map[string]*leveledHook
	level     logrus.Level
	remaining string
	errors    errors.Array // Used to cumultate errors in the logging process
}

type leveledHook struct {
	level logrus.Level
	hook  *Hook
}

func createInnerLogger(reportCaller bool, data logrus.Fields) *logrus.Entry {
	logger := logrus.New()
	logger.ReportCaller = reportCaller
	logger.Out = ioutil.Discard      // Discard all logs to the main logger
	logger.Level = DisabledLevel - 1 // Always log at the highest possible level. Hooks will decide if the log goes through or not
	return logger.WithFields(data)
}

// New creates a new Multilogger instance.
// If no hook is provided, it defaults to standard console logger at warning log level.
func New(module string, hooks ...*Hook) *Logger {
	if len(hooks) == 0 {
		hooks = []*Hook{NewConsoleHook("", logrus.WarnLevel)}
	}
	logger := &Logger{
		Entry:   createInnerLogger(ParseBool(os.Getenv(CallerEnvVar)), logrus.Fields{moduleFieldName: module}),
		Catcher: true,
	}
	logger.AddHooks(hooks...)
	logger.PrintLevel = outputLevel
	return logger
}

// Copy returns a new logger with the same hooks but a different module name.
// module is optional, if not supplied, the original module name will copied.
// If many name are supplied, they are joined with a - separator.
func (logger *Logger) Copy(module ...string) *Logger {
	moduleName := strings.Join(module, "-")
	if len(module) == 0 {
		// The function has been called without argument, so we copy the original module name
		moduleName = logger.GetModule()
	}

	var hooks []*Hook
	for key, hook := range logger.hooks {
		inner := hook.hook.inner
		if cloneable, ok := hook.hook.inner.(cloneable); ok {
			// We duplicate the inner hook if it is cloneable
			inner = cloneable.clone()
		}
		hooks = append(hooks, NewHook(key, hook.level, inner))
	}

	return (&Logger{
		Entry:      createInnerLogger(logger.Logger.ReportCaller, logger.Entry.Data).WithTime(logger.Time).WithContext(logger.Context).WithField(moduleFieldName, moduleName),
		PrintLevel: logger.PrintLevel,
		Catcher:    logger.Catcher,
		level:      logger.level,
		remaining:  logger.remaining,
		errors:     logger.errors,
	}).AddHooks(hooks...)
}

// Child clones the logger, appending the child's name to the parent's module name.
func (logger *Logger) Child(child string) *Logger {
	if module := logger.GetModule(); module != "" {
		return logger.Copy(fmt.Sprintf("%s:%s", module, child))
	}
	return logger.Copy(child)
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

// SetModule sets the module name associated to the current logger.
func (logger *Logger) SetModule(module string) *Logger {
	logger.Data[moduleFieldName] = module
	return logger
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

// TryAddHook adds a hook to the hook collection and associated it with a name and a level.
// Can also be used to replace an existing hook.
func (logger *Logger) TryAddHook(name string, level interface{}, hook logrus.Hook) (*Logger, error) {
	level, err := TryParseLogLevel(level)
	if err != nil {
		return nil, err
	}
	return logger.AddHook(name, level, hook), nil
}

// AddHooks adds a collection of hook wrapper as hook to the current logger.
func (logger *Logger) AddHooks(hooks ...*Hook) *Logger {
	if logger.hooks == nil {
		logger.hooks = make(map[string]*leveledHook)
	}
	for _, hook := range hooks {
		logger.hooks[hook.name] = &leveledHook{hook.level, hook}
		if sl, ok := hook.inner.(setLoggerI); ok {
			// If the hook supports to attach the current logger to it, we set it
			sl.SetLogger(logger)
		}
	}
	return logger.refreshLoggers()
}

// AddConsole adds a console hook to the current logger.
func (logger *Logger) AddConsole(name string, level interface{}, format ...interface{}) *Logger {
	return logger.AddHooks(NewConsoleHook(name, level, format...))
}

// AddFile adds a file hook to the current logger.
func (logger *Logger) AddFile(filename string, isDir bool, level interface{}, format ...interface{}) *Logger {
	return logger.AddHooks(NewFileHook(filename, isDir, level, format...))
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
func (logger *Logger) GetHookLevel(name string) logrus.Level {
	if hook := logger.getHook(name); hook != nil {
		return hook.level
	}
	return DisabledLevel
}

// SetHookLevel set a new log level for a registered hook.
func (logger *Logger) SetHookLevel(name string, level interface{}) error {
	if hook := logger.Hook(name); hook != nil {
		_, err := logger.TryAddHook(hook.name, level, hook.inner)
		return err
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

	for _, key := range logger.ListHooks() {
		hook := logger.hooks[key]
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

// AddError let functions to add error to the current logger indicating that something went
// wrong in the logging process.
func (logger *Logger) AddError(err error) {
	if err != nil {
		logger.errors = append(logger.errors, err)
	}
}

// GetError returns the current error state of the logging process.
func (logger *Logger) GetError() error { return logger.errors.AsError() }

// ClearError cleans up the current error state of the logging process.
// It also returns the current error state.
func (logger *Logger) ClearError() error {
	current := logger.errors.AsError()
	logger.errors = nil
	return current
}

// Close implements io.Closer
func (logger *Logger) Close() error {
	if logger.remaining != "" {
		_, err := logger.Write(nil)
		return err
	}
	return nil
}

// This methods intercepts every message written to stream if Catcher is set and determines if a logging
// function should be used.
func (logger *Logger) Write(writeBuffer []byte) (int, error) {
	if !logger.Catcher {
		if err := logger.printLines(string(writeBuffer)); err != nil {
			return 0, err
		}
		return len(writeBuffer), nil
	}

	var (
		buffer      string
		resultCount int
	)

	if logger.remaining != "" {
		resultCount -= len(logger.remaining)
		buffer = logger.remaining + string(writeBuffer)
		logger.remaining = ""
	} else {
		buffer = string(writeBuffer)
	}

	if writeBuffer != nil {
		lastCR := strings.LastIndex(buffer, "\n")
		logger.remaining = buffer[lastCR+1:]
		buffer = buffer[:lastCR+1]
		resultCount += len(logger.remaining)
	}

	for {
		searchBuffer, extraChar := buffer, 0
		if writeBuffer == nil {
			searchBuffer += "\n"
			extraChar = 1
		}
		matches, _ := reutils.MultiMatch(searchBuffer, logMessages...)
		if len(matches) == 0 {
			break
		}

		if before := matches["before"]; before != "" {
			if err := logger.printLines(before); err != nil {
				return 0, err
			}
			count := len(before)
			resultCount += count
			buffer = buffer[count:]
		}

		level := ParseLogLevel(matches["level"])
		message := matches["message"]
		if prefix := matches["prefix"]; prefix != "" {
			message = fmt.Sprintf("%s %s %s", prefix, level, message)
		}
		logger.Log(level, message)
		if err := logger.GetError(); err != nil {
			return 0, err
		}
		toRemove := len(matches["toRemove"]) - extraChar
		buffer = buffer[toRemove:]
		resultCount += toRemove
	}

	if err := logger.printLines(buffer); err != nil {
		return 0, err
	}
	return resultCount + len(buffer), nil
}

func (logger *Logger) printLines(s string) error {
	lines := strings.Split(s, "\n")
	count := len(lines)
	for i, line := range lines {
		if logger.PrintLevel == outputLevel && i != count-1 {
			logger.Println(line)
		} else if i != count-1 || line != "" {
			logger.Print(line)
		}
		if err := logger.GetError(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	choices := fmt.Sprintf(`\[(?P<level>warn|%s)\]`, strings.Join(AcceptedLevels()[1:], "|"))
	expressions := []string{
		// https://regex101.com/r/jhhPLS/2
		`${choices}[[:blank:]]*{\s*${message}\s*}`,
		`[[:blank:]]*(?P<prefix>[^\n]*?)[[:blank:]]*${choices}[[:blank:]]*${message}[[:blank:]]*\n`,
	}

	for _, expr := range expressions {
		expr = fmt.Sprintf(`(?is)(?P<before>.*?)(?P<toRemove>%s)`, expr)
		expr = strings.Replace(expr, "${choices}", choices, -1)
		expr = strings.Replace(expr, "${message}", `(?P<message>.*?)`, -1)
		logMessages = append(logMessages, regexp.MustCompile(expr))
	}
}

var logMessages []*regexp.Regexp
