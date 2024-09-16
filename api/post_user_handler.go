package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

func (s *Server) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	var dto models.CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := s.store.CreateUser(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
