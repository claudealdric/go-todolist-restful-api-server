package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestServer(t *testing.T) {
	t.Run("responds with 200 OK status on the root path", func(t *testing.T) {
		datastore := &mockDataStore{}
		server := NewServer(datastore)

		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandleRoot)
		handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns the stored tasks with GET on `/tasks`", func(t *testing.T) {
		datastore := &mockDataStore{}
		server := NewServer(datastore)

		request, err := http.NewRequest(http.MethodGet, "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandleGetTasks)
		handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		if datastore.getTasksCalls != 1 {
			t.Errorf(
				"did not receive the right number of calls; want %d, got %d",
				datastore.getTasksCalls,
				1,
			)
		}

		got := testutils.GetTasksFromResponse(t, response.Body)
		want := wantedTasks
		testutils.AssertEquals(t, got, want)
	})
}

var wantedTasks = []models.Task{{Title: "Pack clothes"}}

type mockDataStore struct {
	getTasksCalls int
}

func (m *mockDataStore) GetTasks() []models.Task {
	m.getTasksCalls++
	return wantedTasks
}
