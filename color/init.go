package multicolor

import (
	"strings"

	"github.com/fatih/color"
)

func init() {
	colorNames = make(map[string]color.Attribute, BgHiWhite)
	for i := Reset; i < BgHiWhite; i++ {
		name := strings.ToLower(attribute(i).String())
		if strings.HasPrefix(name, "attribute(") {
			continue
		}
		colorNames[name] = color.Attribute(i)
		if strings.HasPrefix(name, "fg") {
			colorNames[name[2:]] = color.Attribute(i)
		}
	}
}

var colorNames map[string]color.Attribute
