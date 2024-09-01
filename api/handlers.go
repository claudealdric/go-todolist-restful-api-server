package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

const jsonContentType = "application/json"

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

func (s *Server) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	tasks, err := s.store.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandlePatchTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	id, _ := strconv.Atoi(r.PathValue("id")) // TODO: handle error
	task := models.Task{Id: id}
	_ = json.NewDecoder(r.Body).Decode(&task) // TODO: handle error
	updatedTask, _ := s.store.UpdateTask(task)
	_ = json.NewEncoder(w).Encode(updatedTask) // TODO: handle error
}

func (s *Server) HandlePostTasks(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
