package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
)

// TODO: rename
func (s *Server) HandleDeleteTaskById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("ID: %q is invalid", r.PathValue("id")),
			http.StatusBadRequest,
		)
	}
	if err := s.store.DeleteTaskById(id); err != nil {
		if errors.Is(err, data.ErrResourceNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
