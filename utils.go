package multilogger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatDuration returns a string to represent the duration.
func FormatDuration(duration time.Duration) string {
	const day = 24 * time.Hour
	const week = 7 * day
	const month = 30 * day
	const year = 365 * day
	duration = duration.Round(time.Microsecond)
	result := ""
	if duration >= time.Hour {
		duration = duration.Round(time.Second / 10)
	}
	if duration > year {
		result = fmt.Sprintf("%dy", duration/year)
		duration = duration % year
	}
	if duration > 45*day {
		result = fmt.Sprintf("%s%dmo", result, duration/month)
		duration = duration % month
	}
	if duration >= 2*week {
		duration = duration.Round(time.Second)
		result = fmt.Sprintf("%s%dw", result, duration/week)
		duration = duration % week
	}
	if duration > day {
		result = fmt.Sprintf("%s%dd", result, duration/day)
		duration = duration % day
	}
	return fmt.Sprintf("%s%s", result, duration)
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
