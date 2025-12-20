package server

import (
	"net/http"
	"strings"
)

func (s *Server) hasPrefix(path *string) {
	if path == nil {
		return
	}
	if s.prefix != "" {
		if *path == "" {
			*path = s.prefix + *path
			return
		}
		if  strings.HasSuffix(*path, ""){
			*path = strings.TrimSuffix(*path, "/")
		}
		if !strings.HasSuffix(s.prefix, "/") && !strings.HasPrefix(*path, "/") {
			*path = s.prefix + "/" + *path
			return
		}
		if strings.HasSuffix(s.prefix, "/") && !strings.HasPrefix(*path, "/") {
			*path = s.prefix + *path
			return
		}
		if strings.HasSuffix(s.prefix, "/") && strings.HasPrefix(*path, "/") {
			*path = s.prefix + strings.TrimPrefix(*path, "/")
			return
		}
		*path = s.prefix + *path
	} else {
		if !strings.HasPrefix(*path, "/") {
			*path = "/" + *path
		}
	}
}

func (s *Server) Get(path string, handler http.HandlerFunc) {
	s.hasPrefix(&path)
	s.Handler(http.MethodGet+" "+path, handler)
}

func (s *Server) Post(path string, handler http.HandlerFunc) {
	s.hasPrefix(&path)
	s.Handler(http.MethodPost+" "+path, handler)
}

func (s *Server) Put(path string, handler http.HandlerFunc) {
	s.hasPrefix(&path)
	s.Handler(http.MethodPut+" "+path, handler)
}

func (s *Server) Delete(path string, handler http.HandlerFunc) {
	s.hasPrefix(&path)
	s.Handler(http.MethodDelete+" "+path, handler)
}

func (s *Server) Options(path string, handler http.HandlerFunc) {
	s.hasPrefix(&path)
	s.Handler(http.MethodOptions+" "+path, handler)
}
