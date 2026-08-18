package webserve

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

type Server struct {
	addr    int
	mux     *http.ServeMux
	content fs.FS
}

func New(addr int, content fs.FS) *Server {
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

func (s *Server) Get(route string, handler http.HandlerFunc) {
	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}
	s.mux.HandleFunc(fmt.Sprintf("GET %s", route), handler)
}
func (s *Server) Post(route string, handler http.HandlerFunc) {
	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}
	s.mux.HandleFunc(fmt.Sprintf("POST %s", route), handler)
}
func (s *Server) Put(route string, handler http.HandlerFunc) {
	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}
	s.mux.HandleFunc(fmt.Sprintf("PUT %s", route), handler)
}
func (s *Server) Delete(route string, handler http.HandlerFunc) {
	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}
	s.mux.HandleFunc(fmt.Sprintf("DELETE %s", route), handler)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}
