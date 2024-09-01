package api

import (
	"fmt"
	"net/http"
)

func NewRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(get("/{$}"), s.HandleRoot)
	mux.HandleFunc(get("/tasks"), s.HandleGetTasks)
	mux.HandleFunc(get("/tasks/{id}"), s.HandleGetTaskById)
	mux.HandleFunc("PATCH /tasks/{id}", s.HandlePatchTask)
	mux.HandleFunc("POST /tasks", s.HandlePostTasks)
	mux.HandleFunc("DELETE /tasks/{id}", s.HandleDeleteTask)
	return mux
}

func get(path string) string {
	return fmt.Sprintf("%s %s", http.MethodGet, path)
}
