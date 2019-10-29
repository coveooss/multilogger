// +build linux darwin

package multicolor

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/fatih/color"
)

func ExampleSprint() {
	// Sprint is also very flexible.

	// You can use Sprint to format a message in color by prefixing the real text by color attributes.
	fmt.Println(Sprint("Yellow", color.Underline, "Hello"))

	// As with New, you can mix strings & color attribute. You can also use Sprint like a Sprintf function.
	// Color attributes following the message are not considered as attributes and are simply printed.
	fmt.Println(Sprint("yellow+underline", color.BlinkRapid, "color1=%v color2=%s", color.BgBlue, "red"))

	// If the arguments are not compatible with Sprintf, the result shows the error.
	fmt.Println(Sprint("HiRed", "Faint", "Wrong format %d %s", color.BgCyan))

	// As with fmt.Sprint, strings are concatenated and first element following a string is also concatenated.
	fmt.Println(Sprint("Red, CrossedOut", color.BlinkRapid, "Hello", color.BgMagenta, "red", "BOLD", 1, 2, 3))

	// Output:
	// Hello
	// color1=44 color2=red
	// Wrong format %d %s46
	// Hello45redBOLD1 2 3
}

func ExampleSprintln() {
	// Sprint is also very flexible.

	// You can use Sprint to format a message in color by prefixing the real text by color attributes.
	fmt.Print(Sprintln("Yellow", color.Underline, "Hello"))

	// As with New, you can mix strings & color attribute. You can also use Sprint like a Sprintf function.
	// Color attributes following the message are not considered as attributes and are simply printed.
	fmt.Print(Sprintln("yellow+underline", color.BlinkRapid, "color1=%v color2=%s", color.BgBlue, "red"))

	// If the format argument appears to have a format, the function behave like a fmt.Sprintf followed by a fmt.Sprintln.
	fmt.Print(Sprintln("HiRed", "Faint", "Background Hi Magenta is %d", color.BgHiMagenta))

	// If the arguments are not compatible with Sprintf, the result shows the error.
	fmt.Print(Sprintln("HiRed", "Faint", "Wrong format %d %s", color.BgCyan))

	// As with fmt.Sprintln, all elements are separated by a space.
	fmt.Print(Sprintln("Red, CrossedOut", color.BlinkRapid, "Hello", color.BgMagenta, "red", "BOLD", 1, 2, 3))
	// Output:
	// Hello
	// color1=44 color2=red
	// Background Hi Magenta is 105
	// Wrong format %d %s46
	// Hello 45 red BOLD 1 2 3
}

func ExamplePrint() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	// Print is used to output to a color stream.
	Print("Yellow", "Hello")

	// It behaves like regular fmt.Print.
	Print("red", "This", "is", "concatenated", 1, 2, 3)
	// Output:
	// HelloThisisconcatenated1 2 3
}

func ExamplePrintln() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	// Print is used to output to a color stream.
	Println("Yellow", "Hello")

	// It behaves like regular fmt.Println.
	Println("red", "This", "is", "not", "concatenated", 1, 2, 3)
	// Output:
	// Hello
	// This is not concatenated 1 2 3
}

func ExamplePrintf() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	// Printf is used to output to a color stream.
	Printf("Yellow", "Hello")

	// It behaves like regular fmt.Printf.
	// It must have a format string if there is more that one argument after the attributes.
	Printf("red", "Hello %s %d\n", "world", 123)

	// Otherwise, the result is far from beautiful ;-).
	Printf("red", "Hello", "world!", 1, math.Pi)
	// Output:
	// HelloHello world 123
	// Hello%!!(MISSING)(EXTRA string=world!, int=1, float64=3.141592653589793)
}

func ExampleErrorPrint() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetError(os.Stdout)

	// Print is used to output to a color stream.
	ErrorPrint("Yellow", "Hello")

	// It behaves like regular fmt.Print.
	ErrorPrint("red", "This", "is", "concatenated", 1, 2, 3)
	// Output:
	// HelloThisisconcatenated1 2 3
}

func ExampleErrorPrintln() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetError(os.Stdout)

	// ErrorPrintln is used to output to a color stream.
	ErrorPrintln("Yellow", "Hello")

	// It behaves like regular fmt.Println.
	ErrorPrintln("red", "This", "is", "not", "concatenated", 1, 2, 3)
	// Output:
	// Hello
	// This is not concatenated 1 2 3
}

func ExampleErrorPrintf() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetError(os.Stdout)

	// ErrorPrintf is used to output to a color stream.
	ErrorPrintf("Yellow", "Hello")

	// It behaves like regular fmt.Printf.
	// It must have a format string if there is more that one argument after the attributes.
	ErrorPrintf("red", "Hello %s %d\n", "world", 123)

	// Otherwise, the result is far from beautiful ;-).
	ErrorPrintf("red", "Hello", "world!", 1, math.Pi)
	// Output:
	// HelloHello world 123
	// Hello%!!(MISSING)(EXTRA string=world!, int=1, float64=3.141592653589793)
}

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{"No argument", nil, ""},
		{"Empty arguments", []interface{}{}, ""},
		{"Single argument", []interface{}{"Hello"}, "Hello"},
		{"Two arguments", []interface{}{"Hello", "World"}, "Hello World"},
		{"Two arguments with format", []interface{}{"Hello %s! %d", "World", 100}, "Hello World! 100"},
		{"Bad format", []interface{}{"Hello %s! %d", "World"}, "Hello %s! %dWorld"},
		{"Escaped %", []interface{}{"You got %d%% off", 60}, "You got 60% off"},
		{"Escaped % witout format string", []interface{}{"%%", 60}, "%%60"},
		{"Escaped % with error lookalike result", []interface{}{"%%%!", 60}, "%%%!60"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatMessage(tt.args...); got != tt.want {
				t.Errorf("FormatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
