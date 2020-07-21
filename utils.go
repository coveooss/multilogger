package multilogger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

// FormatDurationNative returns a string to represent the duration using the native golang duration format.
func FormatDurationNative(duration time.Duration) string {
	return fmt.Sprintf("%v", duration)
}

// FormatDurationNativeLong returns a string to represent the duration using the native golang duration format using long units.
func FormatDurationNativeLong(duration time.Duration) string {
	return formatDuration(duration, time.Minute, day, true)
}

// FormatDurationPrecise returns a string to represent the duration with precise unit for each segment 1d2h3m4s5ms6µs.
func FormatDurationPrecise(duration time.Duration) string {
	return formatDuration(duration, 0, 0, false)
}

// FormatDurationPreciseLong returns a string to represent the duration with precise unit for each segment using long units.
func FormatDurationPreciseLong(duration time.Duration) string {
	return formatDuration(duration, 0, 0, true)
}

// FormatDurationClassic returns a string to represent the duration using floating representation for values bellow 1 minutes.
func FormatDurationClassic(duration time.Duration) string {
	return formatDuration(duration, time.Minute, 0, false)
}

// FormatDurationClassicLong returns a string to represent the duration using floating representation for values bellow 1 minutes and display long units.
func FormatDurationClassicLong(duration time.Duration) string {
	return formatDuration(duration, time.Minute, 0, true)
}

// FormatDurationRounded returns a string to represent the duration with all units and reduction of the precision for large value.
func FormatDurationRounded(duration time.Duration) string {
	return FormatDurationPrecise(roundedDuration(duration))
}

// FormatDurationRoundedLong returns a string to represent the duration with all units and reduction of the precision for large value and display long units.
func FormatDurationRoundedLong(duration time.Duration) string {
	return FormatDurationPreciseLong(roundedDuration(duration))
}

// FormatDurationRoundedNative returns a string to represent the duration with reduction of the precision for large value.
func FormatDurationRoundedNative(duration time.Duration) string {
	return FormatDurationNative(roundedDuration(duration))
}

// FormatDurationRoundedNativeLong returns a string to represent the duration with reduction of the precision for large value and display long units.
func FormatDurationRoundedNativeLong(duration time.Duration) string {
	return FormatDurationNative(roundedDuration(duration))
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
	case duration >= time.Minute:
		return duration.Round(time.Second)
	case duration >= time.Second:
		return duration.Round(10 * time.Millisecond)
	case duration >= time.Millisecond:
		return duration.Round(time.Millisecond)
	case duration >= time.Microsecond:
		return duration.Round(time.Microsecond)
	}
	return duration
}

// ParseBool returns true for any value that is set and is not clearly a false.
// It never returns an error.
//
// True values are: 1, t, true, yes, y, on (of any non false value)
// False values are: 0, f, false, no, n, off
func ParseBool(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	if result, err := strconv.ParseBool(value); err == nil {
		return result
	}

	switch value {
	case "", "no", "n", "off":
		return false
	}
	return true
}
