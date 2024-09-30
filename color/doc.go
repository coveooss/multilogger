// Copyright 2019 Coveo Solution inc. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package multicolor is a complement to github.com/fatih/color.

# Description

It allows working with colors and attributes by using their string names. It also provides various print functions
(Sprint, Sprintf, Sprintln, Print, Printf, Println, ErrorPrint, ErrorPrintf and Println) to let users provides colors
and attributes before the actual text they want to print.

# Use it as io.Writer

Color objects defined by fatih cannot be used to configure a stream. With multicolor, Color objects can be
configured as out stream or error stream.

	import (
	    "github.com/coveooss/multilogger/color"
	    "github.com/fatih/color"
	)

	var (
	    Println      = multicolor.Println
	    ErrorPrintln = multicolor.ErrorPrintln
	)

	func main() {
	    // You can use a color object to write your output.
	    multicolor.New("Yellow", "BgBlue").SetOut()
	    Println("Hello world!")

	    // You can also convert a Color object from fatih package and use it as a stream.
	    multicolor.NewColorWriter(color.New(color.FgRed)).SetError()
	    ErrorPrintln("There is something wrong!")
	}

# Example

Here is a complete example:

	import (
	    "github.com/coveooss/multilogger/color"
	    "github.com/fatih/color"
	)

	var Println = multicolor.Println

	func main() {
	    // You can use a color object to write your output.
	    c := multicolor.New("Red", "BgWhite")
	    c.Println("Hello world!")

	    // Or you can directly supply your color to the global Print and ErrorPrint functions.
	    // Multiple attributes can be combined, they will be ignored if they are not supported on your OS.
	    multicolor.Println("BgBlue, Underline, CrossedOut, Faint, BlinkSlow, Italic", "The sky is blue.")
	    fmt.Println("Hello, I should be in Hi Green.")

	    // It is also possible to set the colors for all following outputs. The color setting
	    // will then remain until another color is sent to the same output.
	    multicolor.Set("HiGreen")

	    // It is possible to use a mix of strings, multicolor.Attribute or color.Attribute (from fatih)
	    // to configure your color object.
	    // You can configure the default out stream to always print in color.
	    multicolor.New(multicolor.BgGreen, color.FgWhite, "underline").SetOut()
	    Println("I will be printed in color")
	}

The attributes name are not case significant. Any attributes defined in the constants can be used as
string. For simplicity, it is allowed to forget the prefix FG to specify foreground colors.

	Blue      = FgBlue
	HiMagenta = FgHiMagenta
*/
package multicolor
