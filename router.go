package main

import (
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/handlers"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleRoot)
	return mux
}
