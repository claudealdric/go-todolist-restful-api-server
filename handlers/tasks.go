package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	tasks := []models.Task{{Title: "Buy groceries"}}
	json.NewEncoder(w).Encode(tasks)
}
