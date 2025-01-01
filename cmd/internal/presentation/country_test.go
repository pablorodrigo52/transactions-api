package presentation

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateCountry(t *testing.T) {
	tests := []struct {
		name          string
		input         Country
		expectedError *ApiError
	}{
		{name: "Validate Country with success", input: "Brazil", expectedError: nil},
		{name: "Validate Country throws exception", input: "", expectedError: NewApiError(http.StatusBadRequest, "country name is required")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer assertPanicErrors(t, tt.expectedError)

			if tt.expectedError == nil {
				assert.NotPanics(t, func() {
					tt.input.Validate()
				})
			} else {
				assert.Panics(t, func() {
					tt.input.Validate()
				})
			}
		})
	}
}

func Test_NormalizeCountry(t *testing.T) {
	tests := []struct {
		input    Country
		expected string
	}{
		{input: "España", expected: "Espana"},
		{input: "Français", expected: "Francais"},
		{input: "Brazil!", expected: "Brazil!"},
		{input: "México123", expected: "Mexico123"},
		{input: "argentina", expected: "Argentina"},
		{input: "日本", expected: "日本"},
		{input: "日本 日本 日本", expected: "日本 日本 日本"},
		{input: "Burkina Faso", expected: "Burkina Faso"},
		{input: "Áurkina õaso", expected: "Aurkina Oaso"},
	}

	for _, test := range tests {
		response := test.input.Normalize()

		assert.Equal(t, test.expected, response)
	}
}
