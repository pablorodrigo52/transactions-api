package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundPurchaseAmount(t *testing.T) {
	tests := []struct {
		input    float32
		expected float32
	}{
		{input: 1.234, expected: 1.23},
		{input: 1.235, expected: 1.24},
		{input: 1.236, expected: 1.24},
		{input: 1.0, expected: 1.0},
		{input: 0.555, expected: 0.56},
		{input: 0.554, expected: 0.55},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, RoundPurchaseAmount(test.input))
	}
}
