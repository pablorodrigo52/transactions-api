package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONContentTypeMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		target     string
		wantHeader string
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			target:     "/test",
			wantHeader: "application/json",
		},
		{
			name:       "POST request",
			method:     http.MethodPost,
			target:     "/test",
			wantHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := JSONContentTypeMiddleware(nextHandler)

			req := httptest.NewRequest(tt.method, tt.target, nil)
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if got := rr.Header().Get("Content-Type"); got != tt.wantHeader {
				t.Errorf("Content-Type header = %v, want %v", got, tt.wantHeader)
			}
		})
	}
}
