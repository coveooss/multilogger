package multicolor

import "github.com/ghishadow/color"

type (
	// Attribute is a copy of fatih/Attribute used to generate stringable attributes.
	Attribute = color.Attribute
	// Color is imported from ghishadow/color.
	Color     = color.Color
	attribute Attribute
)

// The following constant are copied from the color package in order to get the actual string names.
const (
	Reset attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground attributes
const (
	FgBlack attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground attributes high intensity
const (
	FgHiBlack attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background attributes
const (
	BgBlack attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background attributes high intensity
const (
	BgHiBlack attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

//go:generate stringer -type=attribute -output generated_colors.go
