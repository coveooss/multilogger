package multilogger

import (
	"os"
	"time"
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
	// DurationPrecisionEnvVar defines the duration precision that should be used to render duration
	DurationPrecisionEnvVar = "MULTILOGGER_DURATION_PRECISION"
)

// SetGlobalFormat configure the default format used for console logging and ensure that it is available for all applications
// by setting an environment variable
func SetGlobalFormat(format string, override bool) (string, bool) {
	return setGlobalFormat(FormatEnvVar, format, override)
}

// SetGlobalFileFormat configure the default format used for file logging and ensure that it is available for all applications
// by setting an environment variable
func SetGlobalFileFormat(format string, override bool) (string, bool) {
	return setGlobalFormat(FormatFileEnvVar, format, override)
}

func setGlobalFormat(envvar, format string, override bool) (string, bool) {
	current, isSet := os.LookupEnv(envvar)
	if !isSet || override {
		os.Setenv(envvar, format)
		return format, true
	}
	return current, false
}

// GetGlobalFormat returns the currently globally set console log format
func GetGlobalFormat() string { return os.Getenv(FormatEnvVar) }

// GetGlobalFileFormat returns the currently globally set file log format
func GetGlobalFileFormat() string { return os.Getenv(FormatFileEnvVar) }

// SetDurationPrecision allows duration to be rounded up to the desired precision
func SetDurationPrecision(precision time.Duration, override bool) (time.Duration, bool) {
	if _, isSet := os.LookupEnv(DurationPrecisionEnvVar); override || !isSet {
		roundDuration = precision
		return roundDuration, true
	}
	return roundDuration, false
}

var roundDuration time.Duration

func init() {
	if duration, isSet := os.LookupEnv(DurationPrecisionEnvVar); isSet {
		roundDuration, _ = time.ParseDuration(duration)
	} else {
		roundDuration = time.Millisecond
	}
}
