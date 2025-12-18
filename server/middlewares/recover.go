package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("CRITICAL: Recovered from panic: %v\nStack Trace:\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RecoverMiddlewareJSON(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("PANIC RECOVERED: %v\n%s", err, debug.Stack())

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Internal Server Error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
}
