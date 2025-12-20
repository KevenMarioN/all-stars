/*
Package server provides a lightweight wrapper around the standard net/http library,
specifically designed to enhance the capabilities of http.ServeMux.

It aims to provide a minimal set of abstractions for building RESTful APIs, offering
middleware chaining, simplified route registration, and standardized server configuration,
all while maintaining 100% compatibility with the standard http.Handler interface.

Usage:

	srv := server.New()
	srv.GET("/api/v1/status", statusHandler)
	srv.Run(8080)
*/
package server

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.ServeMux
	globalMW     []func(http.Handler) http.Handler
	groupMW      []func(http.Handler) http.Handler
	isSubGroup   bool
	prefix       string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewServer() *Server {
	return &Server{ServeMux: http.NewServeMux(), globalMW: make([]func(http.Handler) http.Handler, 0)}
}

func (s *Server) WithReadTimeout(r time.Duration) *Server {
	if !s.isSubGroup {
		s.readTimeout = r
	}
	return s
}

func (s *Server) WithWriteTimeout(r time.Duration) *Server {
	if !s.isSubGroup {
		s.readTimeout = r
	}
	return s
}

func (s *Server) Run(port uint) error {
	var (
		readTimeout  = 10 * time.Second
		writeTimeout = 10 * time.Second
	)

	if s.readTimeout > 0 {
		readTimeout = s.readTimeout
	}
	if s.writeTimeout > 0 {
		writeTimeout = s.writeTimeout
	}
	var h http.Handler = s.ServeMux
	for _, mw := range slices.Backward(s.globalMW) {
		h = mw(h)
	}
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      h,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Info().Uint("listen", port).Msg("Server listener")
	return srv.ListenAndServe()
}

func (s *Server) Use(mw ...func(http.Handler) http.Handler) {
	if s.isSubGroup {
		s.groupMW = append(s.groupMW, mw...)
		return
	}
	s.globalMW = append(s.globalMW, mw...)
}

func getFunctionName(handler any) string {
	fullPath := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	parts := strings.Split(fullPath, "/")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, "-fm")
}

func (s *Server) Handler(path string, h http.Handler) {
	router := strings.Split(path, " ")
	if len(router) == 2 {
		funcName := getFunctionName(h)
		log.Info().Msgf("[%s] %s %v", router[0], router[1], funcName)
	}
	for _, mw := range slices.Backward(s.groupMW) {
		h = mw(h)
	}
	s.ServeMux.Handle(path, h)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = s.ServeMux
	for _, mw := range slices.Backward(s.globalMW) {
		h = mw(h)
	}
	h.ServeHTTP(w, r)
}

func (s *Server) Group(prefix string) *Server {
	if prefix != "" {
		if s.prefix == "" {
			if !strings.HasPrefix(prefix, "/") {
				prefix = "/" + prefix
			}
			if !strings.HasSuffix(prefix, "/") {
				prefix = prefix + "/"
			}
		} else {
			if !strings.HasSuffix(s.prefix, "/") {
				s.prefix = s.prefix + "/"
			}
		}
	}

	subgroup := &Server{
		ServeMux:   s.ServeMux,
		groupMW:    slices.Clone(s.groupMW),
		isSubGroup: true,
		prefix:     s.prefix + prefix,
	}
	copy(subgroup.groupMW, s.groupMW)
	return subgroup
}
