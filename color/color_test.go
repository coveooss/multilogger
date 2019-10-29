// +build linux darwin

package multicolor

import (
	"fmt"

	"github.com/fatih/color"
)

func ExampleNew() {
	// It accepts an enumeration of attributes.
	writer := New(color.FgGreen, color.Underline)
	fmt.Printf("%+v", *writer.Color)
	// Output:
	// {params:[32 4] noColor:<nil>}
}

func ExampleNew_with_array() {
	// It also accepts an array of attributes.
	writer := New([]color.Attribute{color.FgGreen, color.Underline})
	fmt.Printf("%+v", *writer.Color)
	// Output:
	// {params:[32 4] noColor:<nil>}
}

func ExampleNew_with_strings() {
	// It also accepts a list of string attributes.
	writer := New("FgHiGreen", "Underline", "BgYellow")
	fmt.Printf("%+v", *writer.Color)
	// Output:
	// {params:[92 4 43] noColor:<nil>}
}

func ExampleNew_mixup() {
	// It can also mix strings and attributes
	writer := New(color.BgRed, "crossedout")
	fmt.Printf("%+v", *writer.Color)
	// Output:
	// {params:[41 9] noColor:<nil>}
}

func ExampleNew_string_separated() {
	// Or separated by any non letter. Case is also non significant.
	writer := New("RED | underline, CrossedOUT+BgYellow")
	fmt.Printf("%+v", *writer.Color)
	// Output:
	// {params:[31 4 9 43] noColor:<nil>}
}

func ExampleNew_error() {
	// But it panics if you supplied invalid attribute.
	func() {
		defer func() { fmt.Println(recover()) }()
		color := New("red, invalid | rouge")
		fmt.Printf("%+v", color)
	}()
	// Output:
	// Attribute not found invalid
	// Attribute not found rouge
}

func ExampleTryNew() {
	// You can avoid panic by using TryNew.
	if _, err := TryNew("FgBLUE Another invalid color BGYellow"); err != nil {
		fmt.Println(err)
	}
	// Output:
	// Attribute not found Another
	// Attribute not found invalid
	// Attribute not found color
}

func ExampleSet() {
	// It accepts an enumeration of attributes.
	Set("Red+BgYellow", color.Underline)
	fmt.Println("I should be colored")

	// But it panics if you supplied invalid attribute.
	func() {
		defer func() { fmt.Println(recover()) }()
		Set("red, invalid | rouge")
	}()
	// Output:
	// I should be colored
	// Attribute not found invalid
	// Attribute not found rouge
}

func ExampleTrySet() {
	// You can avoid panic by using TrySet.
	if _, err := TrySet("FgBLUE Another invalid color BGYellow"); err != nil {
		fmt.Println(err)
	}
	// Output:
	// Attribute not found Another
	// Attribute not found invalid
	// Attribute not found color
}
