package presentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeCountry(t *testing.T) {
	tests := []struct {
		input    Country
		expected string
	}{
		{input: "España", expected: "Espana"},
		{input: "Français", expected: "Francais"},
		{input: "Brazil!", expected: "Brazil"},
		{input: "México123", expected: "Mexico"},
		{input: "日本", expected: "日本"},
		{input: "Burkina Faso", expected: "Burkina Faso"},
	}

	for _, test := range tests {
		response := test.input.Normalize()

		assert.Equal(t, test.expected, response)
	}
}
