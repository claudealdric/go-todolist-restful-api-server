package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	tasks := s.store.GetTasks()
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
