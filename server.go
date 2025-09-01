package kap

import (
	"fmt"
	"net/http"
)

type Server struct {
	Options *Options
}

func NewServer(options *Options) *Server {
	return &Server{
		Options: options,
	}
}

func (server *Server) Run() error {
	http.HandleFunc("/_ping", HandlePing)
	http.Handle("/", NewAuthHandler(server.Options))
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Options.Port), nil)
}
