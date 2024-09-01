package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

func (s *Server) HandlePostTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err = s.store.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
