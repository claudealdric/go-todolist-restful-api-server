package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
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

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	tasks := s.store.GetTasks()
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *Server) HandlePostTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task = s.store.CreateTask(task)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
