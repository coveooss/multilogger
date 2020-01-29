package multilogger

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	multicolor "github.com/coveooss/multilogger/color"
	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultTimestampFormat = "2006/01/02 15:04:05.000"
	defaultLogFormat       = "[%.4level:color,upper%]: %time% - %message%"
	timeEnvVar             = "MULTILOGGER_BASETIME"
)

type formatterI interface {
	SetLogFormat(...string) *Formatter
	SetColor(bool)
}

// NewFormatter creates a new formatter with color setting and takes the first defined format string as the log format.
func NewFormatter(color bool, formats ...interface{}) *Formatter {
	f := &Formatter{color: color}
	f.initOnce.Do(f.init)
	return f.SetLogFormat(formats...)
}

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	// Available standard keys: time, delay, globaldelay, delta, message, level, module, file, line, func.
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %message%
	TimestampFormat string
	FormatDuration  func(time.Duration) string
	FormatCaller    func(*runtime.Frame) string

	// ColorMap allows user to define the color attributes associated with the error level.
	// Attribute names are defined by http://github.com/fatih/color
	ColorMap map[logrus.Level][]multicolor.Attribute

	// LevelName allows user to rename default level name.
	LevelName map[logrus.Level]string

	format         string
	replacer       *replacer
	color          bool
	initOnce       sync.Once
	replacerLock   sync.Mutex
	baseTime, last time.Time
}

// SetLogFormat initialize the log format with the first defined format in the list.
func (f *Formatter) SetLogFormat(formats ...interface{}) *Formatter {
	// We delete the current formatter code replacer
	f.replacer = nil

	// We set the first non empty format as the current format string
	for _, format := range formats {
		if format != "" {
			f.format = fmt.Sprint(format)
			return f
		}
	}

	// If there is no format specified, we use the default log format
	f.format = defaultLogFormat
	return f
}

// SetColor set color mode on the formatter.
func (f *Formatter) SetColor(color bool) { f.color = color }

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.initOnce.Do(f.init)

	output, err := f.doFormat(entry)
	f.last = entry.Time
	return []byte(output), err
}

func (f *Formatter) init() {
	if f.TimestampFormat == "" {
		f.TimestampFormat = defaultTimestampFormat
	}
	if f.format == "" {
		f.format = defaultLogFormat
	}
	if f.FormatDuration == nil {
		f.FormatDuration = FormatDuration
	}
	if f.FormatCaller == nil {
		f.FormatCaller = func(frame *runtime.Frame) (result string) {
			if frame == nil {
				return
			}
			result = frame.Function
			if frame.File != "" {
				if result != "" {
					result += " "
				}
				result += fmt.Sprintf("%s:%d", frame.File, frame.Line)
			}
			return
		}
	}

	if f.ColorMap == nil {
		f.ColorMap = map[logrus.Level][]multicolor.Attribute{
			logrus.PanicLevel: multicolor.Attributes("Magenta", "Bold"),
			logrus.FatalLevel: multicolor.Attributes("Red", "Bold"),
			logrus.ErrorLevel: multicolor.Attributes("Red"),
			logrus.WarnLevel:  multicolor.Attributes("Yellow"),
			logrus.InfoLevel:  multicolor.Attributes("Blue", "Bold"),
			logrus.DebugLevel: multicolor.Attributes("Green"),
			logrus.TraceLevel: multicolor.Attributes("Green", "Faint"),
		}
	}

	f.baseTime = time.Now()
	f.last = f.baseTime
	if globalTime.IsZero() {
		if t := os.Getenv(timeEnvVar); t != "" {
			var err error
			if globalTime, err = time.Parse(time.RFC3339, t); err == nil {
				return
			}
		}
		globalTime = f.baseTime
		os.Setenv(timeEnvVar, globalTime.Format(time.RFC3339))
	}
}

var globalTime time.Time
