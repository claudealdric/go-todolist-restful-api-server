package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestHandleRoot(t *testing.T) {
	t.Run("responds with 200 OK status", func(t *testing.T) {
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
}

func TestHandleDeleteTaskById(t *testing.T) {
	t.Run("responds with 204 No Content", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		request, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/tasks/%d", initialTasks[0].Id),
			nil,
		)
		testutils.AssertNoError(t, err)

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandleDeleteTaskById)
		handler.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusNoContent)
		// testutils.AssertCalls(t, datastore.createTaskCalls, 1)
	})
}

func TestHandleGetTasks(t *testing.T) {
	t.Run("returns the stored tasks", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		request, err := http.NewRequest(http.MethodGet, "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(server.HandleGetTasks)
		handler.ServeHTTP(response, request)

		testutils.AssertContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		testutils.AssertCalls(t, datastore.getTasksCalls, 1)
		testutils.AssertEquals(
			t,
			testutils.GetTasksFromResponse(t, response.Body),
			initialTasks,
		)
	})
}

func TestHandlePostTasks(t *testing.T) {
	t.Run("creates and returns the task with a 201 Status Created", func(t *testing.T) {
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

		testutils.AssertContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		testutils.AssertStatus(t, response.Code, http.StatusCreated)
		testutils.AssertCalls(t, datastore.createTaskCalls, 1)
		testutils.AssertEquals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			newTask,
		)
	})

	t.Run("responds with a 400 Bad Request given an invalid body", func(t *testing.T) {
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
		testutils.AssertCalls(t, datastore.createTaskCalls, 0)
	})
}

var initialTasks = []models.Task{{Title: "Pack clothes"}}

type mockDataStore struct {
	createTaskCalls int
	getTasksCalls   int
	tasks           []models.Task
}

func newMockDataStore() *mockDataStore {
	m := &mockDataStore{tasks: initialTasks}
	return m
}

func (m *mockDataStore) CreateTask(task models.Task) models.Task {
	m.createTaskCalls++
	return task
}

func (m *mockDataStore) GetTasks() []models.Task {
	m.getTasksCalls++
	return m.tasks
}
