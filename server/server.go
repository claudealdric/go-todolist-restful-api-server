package server

import (
	"encoding/json"
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

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	tasks := s.store.GetTasks()
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
