package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[LOGGER] %s %s Started", r.Method, r.URL.Path)

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Printf(
			"[LOGGER] %s %s Finished (%d ms)",
			r.Method,
			r.URL.Path,
			duration.Milliseconds(),
		)
	})
}
