package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/pablorodrigo52/transaction-api/cmd/internal/presentation"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var apiErr *presentation.ApiError
				switch e := err.(type) {
				case *presentation.ApiError:
					apiErr = e
				case error:
					apiErr = presentation.NewApiError(http.StatusInternalServerError, e.Error())
				default:
					apiErr = presentation.NewApiError(http.StatusInternalServerError, "Unknown error")
				}
				w.WriteHeader(apiErr.Code)
				json.NewEncoder(w).Encode(apiErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
