// +build linux darwin

package multicolor

import (
	"fmt"

	"github.com/fatih/color"
)

func ExampleNew() {
	// New is very flexible.

	// It accepts an enumeration of attributes.
	New(color.FgGreen, color.Underline)

	// Or an array of attributes.
	var attributes []color.Attribute = Attributes(color.FgGreen, color.Underline)
	New(attributes)

	// As well as an enumeration of strings.
	New("Green", "Underline")

	// It can also mix strings and attributes
	New(color.BgRed, "crossedout")

	// Attributes can also be separated by spaces.
	New("Green Underline")

	// Or separated by any non letter. Case is also non significant.
	New("RED | underline, CrossedOUT+BgYellow")

	// But it panics if you supplied invalid attribute.
	func() {
		defer func() { fmt.Println(recover()) }()
		New("red, invalid | rouge")
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
