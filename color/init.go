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
		switch name {
		case "reversevideo":
			colorNames["reverse"] = color.Attribute(i)
		case "blinkslow":
			colorNames["blink"] = color.Attribute(i)
		case "concealed":
			colorNames["secret"] = color.Attribute(i)
		case "crossedout":
			colorNames["strikethrough"] = color.Attribute(i)
			colorNames["strike"] = color.Attribute(i)
		}
	}
}

var colorNames map[string]color.Attribute
