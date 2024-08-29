package server

import (
	"fmt"
	"net/http"
)

func NewRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(get("/{$}"), s.HandleRoot)
	mux.HandleFunc(get("/tasks"), s.HandleGetTasks)
	mux.HandleFunc("POST /tasks", s.HandlePostTasks)
	return mux
}

func get(path string) string {
	return fmt.Sprintf("%s %s", http.MethodGet, path)
}
