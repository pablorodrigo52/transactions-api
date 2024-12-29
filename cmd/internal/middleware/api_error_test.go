package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
		expectedError  *presentation.ApiError
	}{
		{
			name: "Error handler without error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name: "Error handler with api error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(presentation.NewApiError(http.StatusBadRequest, "bad request"))
			}),
			expectedStatus: http.StatusBadRequest,
			expectedError:  presentation.NewApiError(http.StatusBadRequest, "bad request"),
		},
		{
			name: "Error handler with generic error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			}),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  presentation.NewApiError(http.StatusInternalServerError, "Unknown error"),
		},
		{
			name: "Error handler with unknown error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(12345)
			}),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  presentation.NewApiError(http.StatusInternalServerError, "Unknown error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := ErrorHandler(tt.handler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != nil {
				var apiErr presentation.ApiError
				err := json.NewDecoder(rr.Body).Decode(&apiErr)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError.Code, apiErr.Code)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			}
		})
	}
}
