package multicolor

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// SetOut lets the user redefine the default writer used to print to stdout.
// This writer will then be used by functions multicolor.Print, multicolor.Printf and multicolor.Println.
//
// The function returns the result color object if it is a colorable stream, otherwise, it returns nil.
func SetOut(out io.Writer) *ColorWriter { return setOutput(out, &color.Output) }

// SetError lets the user redefine the default writer used to print to stderr.
// This writer will then be used by functions multicolor.ErrorPrint, multicolor.ErrorPrintf and multicolor.ErrorPrintln.
//
// The function returns the result color object if it is a colorable stream, otherwise, it returns nil.
func SetError(out io.Writer) *ColorWriter { return setOutput(out, &color.Error) }

// NewColorWriter creates a writeable color from a fatih Color object.
func NewColorWriter(c *color.Color) *ColorWriter {
	return &ColorWriter{c, color.Output}
}

// ColorWriter is used as a regular color.Color object, but it is able to be used as io.Writer.
type ColorWriter struct {
	*color.Color
	out io.Writer
}

// Writer is the implementation of io.Writer. You should not call directly this function.
// The function will fail if called directly on a stream that have not been configured as out stream.
func (c *ColorWriter) Write(p []byte) (n int, err error) {
	if c.out == nil {
		return 0, fmt.Errorf("Color is not configured as writer, call SetOut or SetError before using it")
	}
	return c.out.Write([]byte(c.Sprint(string(p))))
}

// SetOut uses the current color as the default stdout for multicolor.Print functions.
func (c *ColorWriter) SetOut() *ColorWriter {
	return SetOut(c)
}

// SetError uses the current color as the default stderr for multicolor.ErrorPrint functions.
func (c *ColorWriter) SetError() *ColorWriter {
	return SetError(c)
}

func setOutput(out io.Writer, stream *io.Writer) *ColorWriter {
	if c, isColor := out.(interface{}).(*ColorWriter); isColor {
		writer := &ColorWriter{c.Color, *stream}
		*stream = writer
		return writer
	}
	*stream = out
	return nil
}
