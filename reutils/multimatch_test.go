package reutils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiMatch(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name        string
		s           string
		expressions []*regexp.Regexp
		want        map[string]string
		match       int
	}{
		{"empty", "", nil, nil, -1},
		{
			"single", "Hello", []*regexp.Regexp{
				regexp.MustCompile(`H`),
			}, map[string]string{
				"": "H",
			}, 0,
		},
		{
			"hello", "Hello world!", []*regexp.Regexp{
				regexp.MustCompile(`\w+`),
			}, map[string]string{
				"": "Hello",
			}, 0,
		},
		{
			"hello", "Hello world!", []*regexp.Regexp{
				regexp.MustCompile(`^\w+$`),
				regexp.MustCompile(`^Hello (?P<who>\w+).*$`),
			}, map[string]string{
				"":    "Hello world!",
				"who": "world",
			}, 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, match := MultiMatch(tt.s, tt.expressions...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.match, match)
		})
	}
}

func ExampleGetRegexGroup() {
	// Create a group.
	if _, err := NewRegexGroup("test", `^\w+$`, `^Hello (?P<who>\w+).*$`); err != nil {
		fmt.Println("Invalid group definition", err)
	}

	result, match := MultiMatch("Hello world!", GetRegexGroup("test")...)
	fmt.Printf("Matched group %d\n", match)
	fmt.Println(result)
	// Output:
	// Matched group 1
	// map[:Hello world! who:world]
}

func ExampleDeleteRegexGroup() {
	NewRegexGroup("test", `^\w+$`, `^Hello (?P<who>\w+).*$`)
	result := DeleteRegexGroup("test")
	fmt.Println(result)
	// Output:
	// [^\w+$ ^Hello (?P<who>\w+).*$]
}

func ExampleDeleteRegexGroup_not_existing() {
	result := DeleteRegexGroup("not existing")
	fmt.Println(result)
	// Output:
	// []
}

func ExampleNewRegexGroup_with_error() {
	// All supplied regular expressions must be valid.
	if _, err := NewRegexGroup("test", `^\w+$`, `^Hello (?P(who)\w+).*$`); err != nil {
		fmt.Println("Invalid group definition", err)
	}
	// Output:
	// Invalid group definition error parsing regexp: invalid or unsupported Perl syntax: `(?P`
}
