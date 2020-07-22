package multilogger

import (
	"strconv"
	"strings"
)

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
