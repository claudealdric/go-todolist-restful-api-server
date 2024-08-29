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
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) HandlePostTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task = s.store.CreateTask(task)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
