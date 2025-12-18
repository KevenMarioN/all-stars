package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KevenMarioN/all-stars/server/middlewares"
	"github.com/google/uuid"
)

func TestRequestMiddleware(t *testing.T) {
	basic := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	requestHandler := middlewares.RequestIDMiddleware(basic)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	requestHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d (StatusOK), got %d", http.StatusOK, rr.Code)
	}

	requestID := rr.Header().Get(middlewares.REQUEST_ID_HEADER.String())
	if requestID == "" {
		t.Errorf("Request header there is empty")
	}
	if _, err := uuid.Parse(requestID); err != nil {
		t.Errorf("Expected header key %s is valid uuid, got %s", middlewares.REQUEST_ID_HEADER.String(), requestID)
	}
}
