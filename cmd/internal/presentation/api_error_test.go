package presentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApiError(t *testing.T) {
	tests := []struct {
		code     int
		message  string
		expected *ApiError
	}{
		{400, "Bad Request", &ApiError{400, "Bad Request"}},
		{404, "Not Found", &ApiError{404, "Not Found"}},
		{500, "Internal Server Error", &ApiError{500, "Internal Server Error"}},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := NewApiError(tt.code, tt.message)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expected.Error(), result.Error())
		})
	}
}
