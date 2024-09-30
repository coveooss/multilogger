package multilogger

import (
	"fmt"
	"time"

	"github.com/coveooss/multilogger/errors"
)

const (
	day   = 24 * time.Hour
	week  = 7 * day
	month = 30 * day
	year  = 12 * month
)

var units = []struct {
	delay     time.Duration
	divider   time.Duration
	unitShort string
	unitLong  string
}{
	{year, year, "y", "year"},
	{45 * day, month, "mo", "month"},
	{10 * day, week, "w", "week"},
	{day, day, "d", "day"},
	{time.Hour, time.Hour, "h", "hour"},
	{time.Minute, time.Minute, "m", "minute"},
	{time.Second, time.Second, "s", "second"},
	{time.Millisecond, time.Millisecond, "ms", "millisecond"},
	{time.Microsecond, time.Microsecond, "µs", "microsecond"},
	{time.Nanosecond, time.Nanosecond, "ns", "nanosecond"},
}

// DurationFormat represents the representation of a duration rendering format
type DurationFormat int

const (
	// NativeFormat formats durations using the native golang duration format.
	NativeFormat DurationFormat = iota
	// PreciseFormat formats durations with precise unit for each segment 1d2h3m4s5ms6µs.
	PreciseFormat
	// ClassicFormat formats durations with higher units such as days, months, years,
	// but use floating value bellow minutes as native golang format does.
	ClassicFormat
)

// DurationFunc is the prototype used to represent a duration format function.
type DurationFunc func(time.Duration) string

// TryGetDurationFunc returns a function that can be used to format a duration.
func TryGetDurationFunc(format DurationFormat, rounded, longUnit bool) (DurationFunc, error) {
	var (
		minUnit, maxUnit time.Duration
	)
	switch format {
	case PreciseFormat:
	case ClassicFormat:
		minUnit = time.Minute
	case NativeFormat:
		minUnit = time.Minute
		maxUnit = day
	default:
		return func(d time.Duration) string { return fmt.Sprintf("%v", d) }, fmt.Errorf("unknown format %v", format)
	}

	return func(d time.Duration) string {
		if rounded {
			d = roundedDuration(d)
		}
		return formatDuration(d, minUnit, maxUnit, longUnit)
	}, nil
}

// GetDurationFunc returns a function that can be used to format a duration.
func GetDurationFunc(format DurationFormat, rounded, longUnit bool) DurationFunc {
	return errors.Must(TryGetDurationFunc(format, rounded, longUnit)).(DurationFunc)
}

func formatDuration(duration time.Duration, minUnit, maxUnit time.Duration, longUnit bool) (result string) {
	for _, u := range units {
		if u.delay >= maxUnit && maxUnit != 0 {
			continue
		}
		if duration >= u.delay {
			var value float64
			if duration > minUnit {
				div := duration / u.divider
				value = float64(div)
				duration -= div * u.divider
			} else {
				value = float64(duration) / float64(u.divider)
				duration = 0
			}
			if longUnit {
				if result != "" {
					result += " "
				}
				unit := u.unitLong
				if value >= 2 {
					unit += "s"
				}
				result = fmt.Sprintf("%s%v %s", result, float64(value), unit)
			} else {
				result = fmt.Sprintf("%s%v%s", result, float64(value), u.unitShort)
			}
		}
	}
	if result == "" {
		if longUnit {
			result = "0 second"
		} else {
			result = "0s"
		}
	}
	return
}

func roundedDuration(duration time.Duration) time.Duration {
	switch {
	case duration >= month:
		return duration.Round(day)
	case duration >= time.Hour:
		return duration.Round(time.Minute)
	case duration >= 5*time.Minute:
		return duration.Round(time.Second)
	case duration >= time.Minute:
		return duration.Round(5 * time.Second)
	case duration >= 10*time.Second:
		return duration.Round(100 * time.Millisecond)
	case duration >= time.Second:
		return duration.Round(10 * time.Millisecond)
	case duration >= time.Millisecond:
		return duration.Round(10 * time.Microsecond)
	case duration >= time.Microsecond:
		return duration.Round(10 * time.Nanosecond)
	}
	return duration
}
