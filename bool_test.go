package multilogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"True", true},
		{"  True  	\n", true},
		{"T", true},
		{"t", true},
		{"TrUe", true},
		{"On", true},
		{"Y", true},
		{"Yes", true},
		{"Whatever", true},
		{"False", false},
		{"FaLsE", false},
		{"f", false},
		{"off", false},
		{"No", false},
		{"n", false},
		{"NO", false},
		{"	 no  		", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBool(tt.name); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
