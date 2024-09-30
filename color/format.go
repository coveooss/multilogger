package multicolor

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// FormatMessage analyses the arguments to determine if Sprintf or Sprintln should be used.
func FormatMessage(args ...interface{}) string { return formatMessage(sprintlnType, true, args...) }

// Sprint returns a string formatted with attributes that are supplied before.
func Sprint(args ...interface{}) string { return sprint(sprintType, args...) }

// Sprintf returns a string formatted with attributes that are supplied before.
func Sprintf(args ...interface{}) string { return sprint(sprintfType, args...) }

// Sprintln returns a string formatted with attributes that are supplied before.
func Sprintln(args ...interface{}) string { return sprint(sprintlnType, args...) }

// Print call standard fmt.Printf function but using the color out stream.
func Print(args ...interface{}) (int, error) { return fmt.Fprint(color.Output, Sprint(args...)) }

// Println call standard fmt.Println function but using the color out stream.
func Println(args ...interface{}) (int, error) { return fmt.Fprint(color.Output, Sprintln(args...)) }

// Printf behave as Printf, it expects to have a format string after the color attributes.
func Printf(args ...interface{}) (int, error) { return fmt.Fprint(color.Output, Sprintf(args...)) }

// ErrorPrint call standard fmt.Printf function but using the color out stream.
func ErrorPrint(args ...interface{}) (int, error) { return fmt.Fprint(color.Error, Sprint(args...)) }

// ErrorPrintln call standard fmt.Println function but using the color out stream.
func ErrorPrintln(args ...interface{}) (int, error) {
	return fmt.Fprint(color.Error, Sprintln(args...))
}

// ErrorPrintf call standard fmt.Printf function but using the color out stream.
func ErrorPrintf(args ...interface{}) (int, error) {
	return fmt.Fprint(color.Error, Sprintf(args...))
}

func sprint(printType printType, args ...interface{}) string {
	var i int
	var color = new(color.Color)
	for i = 0; i < len(args); i++ {
		if attribute, err := TryConvertAttributes(args[i]); err == nil {
			color.Add(attribute...)
			continue
		}
		break
	}

	output := formatMessage(printType, false, args[i:]...)
	if i == 0 {
		return fmt.Sprint(output)
	}
	return color.Sprint(output)
}

func formatMessage(printType printType, trimEOL bool, args ...interface{}) (result string) {
	var sprintFunc func(...interface{}) string
	switch printType {
	case sprintType, sprintfType:
		sprintFunc = fmt.Sprint
	case sprintlnType:
		sprintFunc = fmt.Sprintln
	}
	switch len(args) {
	case 0, 1:
		result = sprintFunc(args...)
	default:
		if format, newArgs := fmt.Sprint(args[0]), args[1:]; printType == sprintfType || strings.Contains(format, "%") {
			if result = fmt.Sprintf(format, newArgs...); printType != sprintfType && (strings.Count(format, "%!") != strings.Count(result, "%!") || strings.Contains(result, "%!!(")) {
				result = fmt.Sprint(args...)
			}
			if printType == sprintlnType {
				result += "\n"
			}
		}
		if result == "" {
			result = sprintFunc(args...)
		}
	}

	if trimEOL {
		result = strings.TrimRight(result, fmt.Sprintln())
	}
	return
}

type printType uint8

const (
	sprintType printType = iota
	sprintlnType
	sprintfType
)
