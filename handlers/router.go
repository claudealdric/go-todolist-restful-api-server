package handlers

import (
	"fmt"
	"net/http"
)

func NewRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(get("/{$}"), HandleRoot)
	mux.HandleFunc(get("/tasks"), s.HandleGetTasks)
	return mux
}

func get(path string) string {
	return fmt.Sprintf("%s %s", http.MethodGet, path)
}
