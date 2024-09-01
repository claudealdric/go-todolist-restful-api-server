package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
)

func (s *Server) HandleGetTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("ID: %q is invalid", r.PathValue("id")),
			http.StatusBadRequest,
		)
		return
	}
	task, err := s.store.GetTaskById(id)
	if err != nil {
		if errors.Is(err, data.ErrResourceNotFound) {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(
			w,
			fmt.Sprintf("ID: %q is invalid", r.PathValue("id")),
			http.StatusBadRequest,
		)
		return
	}
}
