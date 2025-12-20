package server_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/KevenMarioN/all-stars/server"
)

func TestServer(t *testing.T) {
	var callBy string
	server := server.NewServer()
	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callBy += "1"
			next.ServeHTTP(w, r)
		})
	}
	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callBy += "2"
			next.ServeHTTP(w, r)
		})
	}
	m3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callBy += "3"
			next.ServeHTTP(w, r)
		})
	}
	m4 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callBy += "4"
			next.ServeHTTP(w, r)
		})
	}
	m5 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callBy += "5"
			next.ServeHTTP(w, r)
		})
	}

	server.Use(m1)
	server.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	v1 := server.Group("v1")
	v1.Use(m2)
	v1.Put("on", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	books := v1.Group("books")
	books.Use(m4)
	categoriesBooks := books.Group("categories")
	categoriesBooks.Use(m3)
	categoriesBooks.Use(m5)
	categoriesBooks.Post("", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "created with success"}`))
	})
	categoriesBooks.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		idPath := r.PathValue("id")
		id, err := strconv.Atoi(idPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "delete with success %d"}`, id)
	})
	categoriesBooks.Options("", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "options is true")
	})
	categoriesBooks.Get("financial", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `[{"id": 1},{"id": 3}]`)
	})
	categoriesBooks.Get("/development", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `[{"id": 9}]`)
	})

	scenarios := []struct {
		RequestMethod     string
		RequestPath       string
		ExpectedCallBy    string
		ExpectedStatus    int
		ExpectedPathValue string
	}{
		{
			RequestMethod:  http.MethodGet,
			RequestPath:    "/health",
			ExpectedStatus: http.StatusOK,
			ExpectedCallBy: "1",
		},
		{
			RequestMethod:  http.MethodPut,
			RequestPath:    "/v1/on",
			ExpectedStatus: http.StatusNoContent,
			ExpectedCallBy: "12",
		},
		{
			RequestMethod:     http.MethodPost,
			RequestPath:       "/v1/books/categories",
			ExpectedStatus:    http.StatusCreated,
			ExpectedCallBy:    "12435",
			ExpectedPathValue: `{"message": "created with success"}`,
		},
		{
			RequestMethod:     http.MethodDelete,
			RequestPath:       "/v1/books/categories/98",
			ExpectedStatus:    http.StatusOK,
			ExpectedCallBy:    "12435",
			ExpectedPathValue: `{"message": "delete with success 98"}`,
		},
		{
			RequestMethod:     http.MethodOptions,
			RequestPath:       "/v1/books/categories",
			ExpectedStatus:    http.StatusAccepted,
			ExpectedCallBy:    "12435",
			ExpectedPathValue: "options is true",
		},
		{
			RequestMethod:     http.MethodGet,
			RequestPath:       "/v1/books/categories/financial",
			ExpectedStatus:    http.StatusOK,
			ExpectedCallBy:    "12435",
			ExpectedPathValue: `[{"id": 1},{"id": 3}]`,
		},
		{
			RequestMethod:     http.MethodGet,
			RequestPath:       "/v1/books/categories/development",
			ExpectedStatus:    http.StatusOK,
			ExpectedCallBy:    "12435",
			ExpectedPathValue: `[{"id": 9}]`,
		},
	}
	for _, tt := range scenarios {
		callBy = ""
		rq, err := http.NewRequest(tt.RequestMethod, tt.RequestPath, nil)
		if err != nil {
			t.Errorf("NewRequest: %s", err.Error())
		}
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, rq)
		rs := rr.Result()
		if rs.StatusCode != tt.ExpectedStatus {
			t.Errorf("[%s] %s: status: expected %d; got %d\n",
				tt.RequestMethod,
				tt.RequestPath,
				tt.ExpectedStatus, rs.StatusCode)
		}
		if callBy != tt.ExpectedCallBy {
			t.Errorf("[%s] %s: mw used: expected %q; got %q\n",
				tt.RequestMethod,
				tt.RequestPath,
				tt.ExpectedCallBy, callBy)
		}
		if tt.ExpectedPathValue != "" {
			rawAll, _ := io.ReadAll(rs.Body)
			all := string(rawAll)
			if all != tt.ExpectedPathValue {
				t.Errorf("[%s] %s: expected path: expected %q; got %q\n",
					tt.RequestMethod,
					tt.RequestPath,
					tt.ExpectedPathValue, all)
			}
		}
	}
}
