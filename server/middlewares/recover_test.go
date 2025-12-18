package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Intentional panic for testing purposes")
	})

	protectedHandler := RecoverMiddleware(panicHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The middleware FAILED to recover from panic: %v", r)
		}
	}()

	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d (Internal Server Error), got %d", http.StatusInternalServerError, rr.Code)
	}

	expectedBody := "Internal Server Error\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}

func TestRecoverMiddlewareJson(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Intentional panic for testing purposes")
	})

	protectedHandler := RecoverMiddlewareJSON(panicHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The middleware FAILED to recover from panic: %v", r)
		}
	}()

	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d (Internal Server Error), got %d", http.StatusInternalServerError, rr.Code)
	}

	expectedBody := `{"error": "Internal Server Error"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}
}
