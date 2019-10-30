package multicolor

import (
	"fmt"
	"strings"
	"unicode"
)

// Attributes convert any object representation to valid color attributes.
// It will panic if an invalid attribute is provided.
func Attributes(attributes ...interface{}) []Attribute {
	result, err := TryConvertAttributes(attributes...)
	if err != nil {
		panic(err)
	}
	return result
}

// TryConvertAttributes tries to convert any object representation to valid color attribute.
// It returns an error if some parameters cannot be converted to valid attributes.
func TryConvertAttributes(attributes ...interface{}) ([]Attribute, error) {
	var errors []string
	result := make([]Attribute, 0, len(attributes))

	idFunc := func(c rune) bool { return !unicode.IsLetter(c) && !unicode.IsNumber(c) }
	for _, attribute := range attributes {
		if attr, ok := attribute.(Attribute); ok {
			attribute = attr
		}
		if attr, ok := attribute.(Attribute); ok {
			result = append(result, attr)
			continue
		}
		if attributes, ok := attribute.([]Attribute); ok {
			result = append(result, attributes...)
			continue
		}
		for _, attribute := range strings.FieldsFunc(fmt.Sprint(attribute), idFunc) {
			if attr, match := colorNames[strings.ToLower(attribute)]; match {
				result = append(result, attr)
			} else {
				errors = append(errors, fmt.Sprintf("Attribute not found %s", attribute))
			}
		}
	}
	if len(errors) > 0 {
		return result, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return result, nil
}
