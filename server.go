package main

import (
	"fmt"
	"io/fs"
	"net/http"
)

type Server struct {
	addr    int8
	mux     *http.ServeMux
	content fs.FS
}

func New(addr int8, content fs.FS) *Server {
	s := &Server{
		addr:    addr,
		mux:     http.NewServeMux(),
		content: content,
	}
	s.registerStaticRoutes()
	return s
}

func (s *Server) registerStaticRoutes() {
	fmt.Println("Serving frontend files from embedded source")
	registerEmbeddedStaticRoutes(s.mux, s.content)
}

func (s *Server) Run() error {
	handler := GzipMiddleware(s.mux)
	return http.ListenAndServe(fmt.Sprintf(":%v", s.addr), handler)
}
