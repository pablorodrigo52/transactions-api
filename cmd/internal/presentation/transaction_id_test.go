package presentation

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TransactionID_Validate(t *testing.T) {
	tests := []struct {
		name          string
		input         TransactionID
		expectedError *ApiError
	}{
		{name: "Validate TransactionID with success", input: "1", expectedError: nil},
		{name: "Validate TransactionID empty", input: "", expectedError: NewApiError(http.StatusBadRequest, "transaction ID is required")},
		{name: "Validate TransactionID not a number", input: "NaN", expectedError: NewApiError(http.StatusBadRequest, "transaction ID must be a valid number: mock error")},
		{name: "Validate TransactionID negative", input: "-1", expectedError: NewApiError(http.StatusBadRequest, "transaction ID must be a valid number")},
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

func Test_TransactionID_Get(t *testing.T) {

	tests := []struct {
		name  string
		input TransactionID
	}{
		{name: "Get TransactionID with success", input: "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := tt.input.Get()

			assert.NotNil(t, response)
			assert.Equal(t, int64(1), response)
		})
	}
}
