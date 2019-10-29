package multicolor

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func ExampleSetOut() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	SetOut(New("BgGreen+Yellow"))
	Println("I should be colored")
	// Output:
	// I should be colored
}

func ExampleSetError() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetError(os.Stdout)

	SetError(New("BgRed+Green"))
	ErrorPrintln("I should be colored")
	// Output:
	// I should be colored
}

func ExampleNewColorWriter() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	c := color.New(color.FgYellow, color.BgHiWhite, color.Italic)
	SetOut(NewColorWriter(c))
	Println("I should be colored")
	// Output:
	// I should be colored
}

func ExampleColorWriter_SetOut() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	// SetOut can be called as a method on a Color object. It can even be used directly
	// on the same line.
	New("BgGreen+Yellow").SetOut().Println("Hello!")
	Println("I should be colored")
	// Output:
	// Hello!
	// I should be colored
}

func ExampleColorWriter_SetError() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)
	SetError(os.Stdout)

	// SetOut can be called as a method on a Color object. It can even be used directly
	// on the same line.
	New("BgGreen+Yellow").SetError().Println("Hello!")
	ErrorPrint("I should be colored")
	Println(", while I shouldn't")
	// Output:
	// Hello!
	// I should be colored, while I shouldn't
}

func ExampleColorWriter_Write() {
	// We set the output because go test framework redirect os.Stdout and ignore os.Stderr
	SetOut(os.Stdout)

	color := New("BgGreen+Yellow")
	n, err := color.Write([]byte("Not set yet\n"))
	fmt.Printf("n=%d, err=%v\n", n, err)

	color = color.SetOut()
	n, err = color.Write([]byte("Should be ok\n"))
	fmt.Printf("n=%d, err=%v\n", n, err)
	// Output:
	// n=0, err=Color is not configured as writer, call SetOut or SetError before using it
	// Should be ok
	// n=13, err=<nil>
}
