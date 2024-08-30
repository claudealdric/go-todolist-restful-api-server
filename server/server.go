package server

import (
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
)

type Server struct {
	store datastore.DataStore
	http.Handler
}

func NewServer(store datastore.DataStore) *Server {
	server := &Server{store: store}
	router := NewRouter(server)
	server.Handler = router
	return server
}
