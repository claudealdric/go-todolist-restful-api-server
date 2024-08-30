package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

const jsonContentType = "application/json"

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

func (s *Server) HandleDeleteTaskById(w http.ResponseWriter, r *http.Request) {
	// TODO: handle error
	id, _ := strconv.Atoi(r.PathValue("id"))
	s.store.DeleteTaskById(id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	tasks, _ := s.store.GetTasks() // TODO: handle error
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	task, _ = s.store.CreateTask(task) // TODO: handle error
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
