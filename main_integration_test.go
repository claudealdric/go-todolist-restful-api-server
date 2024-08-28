package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestServer(t *testing.T) {
	server := NewServer()

	t.Run("responds with a 200 OK status on the root path", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("an unexpected error occurred: %v", err)
		}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("responds with a 404 not found status on the root path", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/not-found", nil)
		if err != nil {
			t.Fatalf("an unexpected error occurred: %v", err)
		}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns a slice of tasks with GET on `/tasks`", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/tasks", nil)
		if err != nil {
			t.Errorf("an error occurred during the request: %v", err)
		}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		got := testutils.GetTasksFromResponse(t, response.Body)
		want := []models.Task{{Title: "Buy groceries"}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func getTasksFromResponse(t *testing.T, body io.Reader) (tasks []models.Task) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&tasks)

	if err != nil {
		t.Fatalf(
			"unable to parse response from server %q into slice of Task: %v",
			body,
			err,
		)
	}

	return tasks
}
