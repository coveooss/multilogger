package multilogger

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var baseTime = time.Date(2018, 6, 24, 12, 34, 56, int(789*time.Millisecond), time.UTC)

func getTestLogger(name string, level ...interface{}) *Logger {
	// We fix output, format & time to get consistent results
	logrus.SetOutput(os.Stdout)
	if level == nil {
		level = []interface{}{logrus.WarnLevel}
	}

	color := true
	// Unset environment variable unless MULTILOGGER_TEST_KEEP_ENV is set to a true value
	keep := os.Getenv("MULTILOGGER_TEST_KEEP_ENV")
	if result, err := strconv.ParseBool(keep); keep == "" || err == nil && !result {
		os.Unsetenv(FormatEnvVar)
		os.Unsetenv(FormatFileEnvVar)
		os.Unsetenv(CallerEnvVar)
		color = false
	}

	return New(name, NewConsoleHook("", level[0]).SetColor(color)).WithTime(baseTime)
}

func ExampleNew_default() {
	log := getTestLogger("default")
	log.Warning("This is a warning")
	log.Println("The logging level is set to", log.GetLevel())
	log.Printf("Module = %s\n", log.GetModule())
	// Output:
	// [default] 2018/06/24 12:34:56.789 WARNING  This is a warning
	// The logging level is set to warning
	// Module = default
}

func ExampleNew_settingLoggingLevel() {
	// Logging level could be set by explicitely declaring the hook
	log := New("console", NewConsoleHook("", logrus.InfoLevel))
	log.Println("The logging level is set to", log.GetLevel())

	// Or it can also be set after initializing the logger
	// It is possible to use either a logrus.Level or a string to specify the level
	log = New("console")
	log.SetHookLevel("", "trace")
	log.Println("The logging level is set to", log.GetLevel())
	// Output:
	// The logging level is set to info
	// The logging level is set to trace
}

func ExampleLogger_Copy() {
	log := getTestLogger("original", logrus.TraceLevel)
	log.Info("Log from original")
	log.Copy("copy").Trace("Log from copy")
	log.Copy("").Debug("I have no module")
	// Output:
	// [original] 2018/06/24 12:34:56.789 INFO     Log from original
	// [copy] 2018/06/24 12:34:56.789 TRACE    Log from copy
	// 2018/06/24 12:34:56.789 DEBUG    I have no module
}

func ExampleLogger_WithTime() {
	log := getTestLogger("time", logrus.InfoLevel)

	// We can create a logger with a fix moment in time.
	t, _ := time.Parse(time.RFC3339, "2020-12-25T00:00:00Z")
	log = log.WithTime(t)
	log.Info("Log from fixed time")
	// Output:
	// [time] 2020/12/25 00:00:00.000 INFO     Log from fixed time
}

func ExampleLogger_AddTime() {
	log := getTestLogger("time", "Trace")

	// We can create a logger with a fix moment in time.
	t, _ := time.Parse(time.RFC3339, "2020-12-25T00:00:00Z")
	log = log.WithTime(t)
	log.Info("Log from fixed time")
	log.AddTime(5 * time.Second).Trace("Log 5 seconds later")
	log.AddTime(8 * time.Millisecond).Warning("Log 8 more milliseconds later")
	// Output:
	// [time] 2020/12/25 00:00:00.000 INFO     Log from fixed time
	// [time] 2020/12/25 00:00:05.000 TRACE    Log 5 seconds later
	// [time] 2020/12/25 00:00:05.008 WARNING  Log 8 more milliseconds later
}

func ExampleLogger_WithField() {
	log := getTestLogger("field", "Trace")

	// We set the format of the log to include fields
	log.Hook("").SetFormat("%module:square% %time% %level:upper% %message% %fields%.")

	// We create a new logger with additional context
	log2 := log.WithField("hello", "world!").WithField("pi", math.Pi)
	log.Info("No additional field")
	log2.Info("With additional fields")
	// Output:
	// [field] 2018/06/24 12:34:56.789 INFO No additional field .
	// [field] 2018/06/24 12:34:56.789 INFO With additional fields hello=world! pi=3.141592653589793.
}

func ExampleLogger_WithFields() {
	log := getTestLogger("field", "Trace")

	// We set the format of the log to include fields
	log.Hook("").SetFormat("%module:square% %time% %level:upper% %message% %fields%.")

	// We create a new logger with additional context
	log2 := log.WithFields(logrus.Fields{
		"hello": "world!",
		"pi":    math.Pi,
	})
	log.Info("No additional field")
	log2.Info("With additional fields")
	// Output:
	// [field] 2018/06/24 12:34:56.789 INFO No additional field .
	// [field] 2018/06/24 12:34:56.789 INFO With additional fields hello=world! pi=3.141592653589793.
}

func ExampleLogger_AddConsole() {
	log := getTestLogger("json")

	// We add an additional console hook.
	log.AddConsole("json", logrus.WarnLevel, new(logrus.JSONFormatter))
	log.Warning("New JSON log")
	// Output:
	// [json] 2018/06/24 12:34:56.789 WARNING  New JSON log
	// {"level":"warning","module-field":"json","msg":"New JSON log","time":"2018-06-24T12:34:56Z"}
}

func ExampleLogger_AddConsole_overwrite() {
	log := getTestLogger("json")

	// We overwrite the default console hook by not specifying a name to our new console.
	// We also set the JSON formatter to pretty format the JSON code.
	log.AddConsole("", logrus.WarnLevel, &logrus.JSONFormatter{PrettyPrint: true})
	log.Warning("New JSON log")
	// Output:
	// {
	//   "level": "warning",
	//   "module-field": "json",
	//   "msg": "New JSON log",
	//   "time": "2018-06-24T12:34:56Z"
	// }
}
