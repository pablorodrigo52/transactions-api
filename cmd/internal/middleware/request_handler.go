package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	log := slog.Default()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path))
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
