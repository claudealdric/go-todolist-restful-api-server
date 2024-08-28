package main

import (
	"fmt"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/handlers"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(get("/{$}"), handlers.HandleRoot)
	mux.HandleFunc(get("/tasks"), handlers.HandleGetTasks)
	return mux
}

func get(path string) string {
	return fmt.Sprintf("%s %s", http.MethodGet, path)
}
