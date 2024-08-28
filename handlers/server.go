package handlers

import (
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
