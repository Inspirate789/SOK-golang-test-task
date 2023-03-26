package server

import (
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(listenAddrApi string, router http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    listenAddrApi,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Close()
}
