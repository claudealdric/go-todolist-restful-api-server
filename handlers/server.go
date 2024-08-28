package handlers

import (
	"fmt"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
)

type Server struct {
	store datastore.DataStore
	http.Handler
}

func NewServer(store datastore.DataStore) *Server {
	server := new(Server)
	server.store = store
	router := NewRouter(server)
	server.Handler = router
	return server
}

func NewRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(get("/{$}"), HandleRoot)
	mux.HandleFunc(get("/tasks"), s.HandleGetTasks)
	return mux
}

func get(path string) string {
	return fmt.Sprintf("%s %s", http.MethodGet, path)
}
