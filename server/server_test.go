package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestServer(t *testing.T) {
	t.Run("responds with 200 OK status on the root path", func(t *testing.T) {
		datastore := newMockDataStore()
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
		datastore := newMockDataStore()
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
		want := initialTasks
		testutils.AssertEquals(t, got, want)
	})

	t.Run("creates and returns the task with a 201 status created with POST on `/tasks`", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		newTask := models.Task{Title: "Exercise"}
		jsonData, err := json.Marshal(newTask)
		testutils.AssertNoError(t, err)
		request, err := http.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer(jsonData),
		)
		testutils.AssertNoError(t, err)

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandlePostTasks)
		handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusCreated)

		if datastore.createTaskCalls != 1 {
			t.Errorf("got %d, want %d", datastore.createTaskCalls, 1)
		}

		got := testutils.GetTaskFromResponse(t, response.Body)
		want := newTask
		testutils.AssertEquals(t, got, want)

	})

	t.Run("responds with a 400 Bad Request given an invalid body with POST on `/tasks`", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		invalidJson := `{`
		request, err := http.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer([]byte(invalidJson)),
		)
		testutils.AssertNoError(t, err)

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandlePostTasks)
		handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)

		if datastore.createTaskCalls != 0 {
			t.Errorf("got %d, want %d", datastore.createTaskCalls, 0)
		}
	})
}

var initialTasks = []models.Task{{Title: "Pack clothes"}}

type mockDataStore struct {
	createTaskCalls int
	getTasksCalls   int
	tasks           []models.Task
}

func newMockDataStore() *mockDataStore {
	m := mockDataStore{}
	m.tasks = initialTasks
	return &m
}

func (m *mockDataStore) CreateTask(task models.Task) models.Task {
	m.createTaskCalls++
	return task
}

func (m *mockDataStore) GetTasks() []models.Task {
	m.getTasksCalls++
	return initialTasks
}
