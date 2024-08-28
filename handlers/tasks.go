package handlers

import "net/http"

func HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
