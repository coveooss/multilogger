package multilogger

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
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
	if !ParseBool(os.Getenv("MULTILOGGER_TEST_KEEP_ENV")) {
		os.Unsetenv(FormatEnvVar)
		os.Unsetenv(FormatFileEnvVar)
		os.Unsetenv(CallerEnvVar)
		color = false
	}

	globalTime = baseTime
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
	log.Copy().Debug("I have the same module as the original")
	// Output:
	// [original] 2018/06/24 12:34:56.789 INFO     Log from original
	// [copy] 2018/06/24 12:34:56.789 TRACE    Log from copy
	// 2018/06/24 12:34:56.789 DEBUG    I have no module
	// [original] 2018/06/24 12:34:56.789 DEBUG    I have the same module as the original
}

func ExampleLogger_Child() {
	log := getTestLogger("original", logrus.TraceLevel)
	log.Info("Log from original")
	log.Child("1").Trace("Log from first child")
	log.Child("2").Trace("Log from second child")
	// Output:
	// [original] 2018/06/24 12:34:56.789 INFO     Log from original
	// [original:1] 2018/06/24 12:34:56.789 TRACE    Log from first child
	// [original:2] 2018/06/24 12:34:56.789 TRACE    Log from second child
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
	log.SetFormat("%module:square% %time% %level:upper% %message% %fields%.")

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
	log.SetFormat("%module:square% %time% %level:upper% %message% %fields%.")

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

func ExampleLogger_AddFile() {
	log := getTestLogger("file")

	var logfile string
	if temp, err := ioutil.TempFile("", "example"); err != nil {
		log.Fatal(err)
	} else {
		logfile = temp.Name()
		defer os.Remove(logfile)
	}

	log.AddFile(logfile, false, logrus.TraceLevel)
	log.Info("This is information")
	log.Warning("This is a warning")

	content, _ := ioutil.ReadFile(logfile)
	fmt.Println("Content of the log file is:")
	fmt.Println(string(content))
	// Output:
	// [file] 2018/06/24 12:34:56.789 WARNING  This is a warning
	// Content of the log file is:
	//
	// # 2018/06/24 12:34:56.789
	// [file] 2018/06/24 12:34:56.789 INFO     This is information
	// [file] 2018/06/24 12:34:56.789 WARNING  This is a warning
}

func ExampleLogger_AddFile_folder() {
	log := getTestLogger("file")

	logDir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(logDir)

	// Adding a log folder and creating a child logger
	log.AddFile(logDir, true, logrus.TraceLevel)
	childLogger := log.Child("folder/module")

	// Logging into the main logger and the child logger
	log.Info("This is information")
	childLogger.Warning("This is a warning")
	childLogger.Info("This is information")

	// Reading the main logger logs
	firstFile := filepath.Join(logDir, "file.log")
	firstContent, _ := ioutil.ReadFile(firstFile)
	fmt.Println("Content of the first log file is:")
	fmt.Println(string(firstContent))

	// Reading the child logger logs
	secondFile := filepath.Join(logDir, "file.folder", "module.log")
	secondContent, _ := ioutil.ReadFile(secondFile)
	fmt.Println("Content of the second log file is:")
	fmt.Println(string(secondContent))
	// Output:
	// [file:folder/module] 2018/06/24 12:34:56.789 WARNING  This is a warning
	// Content of the first log file is:
	// # 2018/06/24 12:34:56.789
	// [file] 2018/06/24 12:34:56.789 INFO     This is information
	//
	// Content of the second log file is:
	// # 2018/06/24 12:34:56.789
	// [file:folder/module] 2018/06/24 12:34:56.789 WARNING  This is a warning
	// [file:folder/module] 2018/06/24 12:34:56.789 INFO     This is information
}

func ExampleLogger_AddFile_folderWithInvalidModuleName() {
	// Create a test logger with lots of special chars in its name
	loggerName := "/abc:def!/g$%?&*().,;`^<>/"
	log := getTestLogger(loggerName)

	logDir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(logDir)

	// Adding a log folder
	log.AddFile(logDir, true, logrus.TraceLevel)

	// Logging into the main logger and the child logger
	log.Info("This is information")

	// Reading the logs (all the special chars except OS separators and module separators (:) will be removed from the file name)
	firstFile := filepath.Join(logDir, "abc.def", "g.log")
	firstContent, _ := ioutil.ReadFile(firstFile)
	fmt.Println(string(firstContent))

	// Output:
	// # 2018/06/24 12:34:56.789
	// [/abc:def!/g$%?&*().,;`^<>/] 2018/06/24 12:34:56.789 INFO     This is information
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

func ExampleLogger_Formatter_roundDuration() {
	log := getTestLogger("field", "Trace")

	// We set the format of the log to include fields
	log.SetFormat("%time% %globaldelay:round% %level:upper% %message%")
	log.Formatter().RoundDuration = 5 * time.Millisecond

	log.Info("Starting")
	for i := time.Duration(1); i < 24*time.Hour; {
		i *= 10
		log.WithTime(baseTime.Add(i)).Infof("%v later", i)
	}

	// Output:
	// 2018/06/24 12:34:56.789 (<5ms) INFO Starting
	// 2018/06/24 12:34:56.789 (<5ms) INFO 10ns later
	// 2018/06/24 12:34:56.789 (<5ms) INFO 100ns later
	// 2018/06/24 12:34:56.789 (<5ms) INFO 1µs later
	// 2018/06/24 12:34:56.789 (<5ms) INFO 10µs later
	// 2018/06/24 12:34:56.789 (<5ms) INFO 100µs later
	// 2018/06/24 12:34:56.790 (<5ms) INFO 1ms later
	// 2018/06/24 12:34:56.799 (10ms) INFO 10ms later
	// 2018/06/24 12:34:56.889 (100ms) INFO 100ms later
	// 2018/06/24 12:34:57.789 (1s) INFO 1s later
	// 2018/06/24 12:35:06.789 (10s) INFO 10s later
	// 2018/06/24 12:36:36.789 (1m40s) INFO 1m40s later
	// 2018/06/24 12:51:36.789 (16m40s) INFO 16m40s later
	// 2018/06/24 15:21:36.789 (2h46m40s) INFO 2h46m40s later
	// 2018/06/25 16:21:36.789 (1d3h46m40s) INFO 27h46m40s later
}

func ExampleSetDurationPrecision() {
	log := getTestLogger("field", "Trace")

	// We set the format of the log to include fields
	log.SetFormat("%time% %globaldelay:round% %level:upper% %message%")
	SetDurationPrecision(2*time.Second, true)

	log.Info("Starting")
	for i := time.Duration(1); i < 24*time.Hour; {
		i *= 22
		log.WithTime(baseTime.Add(i)).Infof("%v later", i)
	}

	// Output:
	// 2018/06/24 12:34:56.789 (<2s) INFO Starting
	// 2018/06/24 12:34:56.789 (<2s) INFO 22ns later
	// 2018/06/24 12:34:56.789 (<2s) INFO 484ns later
	// 2018/06/24 12:34:56.789 (<2s) INFO 10.648µs later
	// 2018/06/24 12:34:56.789 (<2s) INFO 234.256µs later
	// 2018/06/24 12:34:56.794 (<2s) INFO 5.153632ms later
	// 2018/06/24 12:34:56.902 (<2s) INFO 113.379904ms later
	// 2018/06/24 12:34:59.283 (2s) INFO 2.494357888s later
	// 2018/06/24 12:35:51.664 (54s) INFO 54.875873536s later
	// 2018/06/24 12:55:04.058 (20m8s) INFO 20m7.269217792s later
	// 2018/06/24 19:57:36.711 (7h22m40s) INFO 7h22m39.922791424s later
	// 2018/07/01 06:53:35.090 (6d18h18m38s) INFO 162h18m38.301411328s later
}
