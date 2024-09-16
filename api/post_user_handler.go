package api

import (
	"encoding/json"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

func (s *Server) HandlePostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	var dto models.CreateUserDTO
	json.NewDecoder(r.Body).Decode(&dto)
	user, _ := s.store.CreateUser(dto)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
