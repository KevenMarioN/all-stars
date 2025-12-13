package server

import "net/http"

func (s *Server) hasPrefix(path *string) {
	if path == nil {
		return
	}
	if s.prefix != "" {
		*path = s.prefix + *path
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
