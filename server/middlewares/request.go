package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	RequestIDHeader HeaderKey = "X-Request-ID"
	RequestIDKey    MwKey     = "request_id"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		w.Header().Set(RequestIDHeader.String(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
