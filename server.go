package main

import "net/http"

type Server struct {
	http.Handler
}

func NewServer() *Server {
	server := new(Server)
	router := NewRouter()
	server.Handler = router
	return server
}
