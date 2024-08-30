package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestHandleRoot(t *testing.T) {
	t.Run("responds with 200 OK status", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})
}

func TestHandleDeleteTaskById(t *testing.T) {
	t.Run("deletes the task and responds with 204 No Content", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		tasks, err := datastore.GetTasks()
		testutils.AssertNoError(t, err)
		initialTasksCount := len(tasks)
		if initialTasksCount == 0 {
			t.Error("expected at least one initial task")
		}

		taskToDelete := tasks[0]
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/tasks/%d", taskToDelete.Id),
			nil,
		)
		response := httptest.NewRecorder()

		server.Handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNoContent)

		tasks, err = datastore.GetTasks()
		testutils.AssertNoError(t, err)
		if len(tasks) != initialTasksCount-1 {
			t.Errorf(
				"expected a slice of length %d; received %+v",
				initialTasksCount-1,
				tasks,
			)
		}
	})

	// TODO: implement
	// t.Run("responds with 404 Not Found when the task does not exist", func(t *testing.T) {
	// datastore := newMockDataStore()
	// server := NewServer(datastore)
	//
	// taskToDelete := models.Task{100, "Does not exist"}
	// request := httptest.NewRequest(
	// 	http.MethodDelete,
	// 	fmt.Sprintf("/tasks/%d", taskToDelete.Id),
	// 	nil,
	// )
	// response := httptest.NewRecorder()
	//
	// server.Handler.ServeHTTP(response, request)
	// testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	// })
}

func TestHandleGetTasks(t *testing.T) {
	t.Run("returns the stored tasks", func(t *testing.T) {
		datastore := newMockDataStore()
		server := NewServer(datastore)

		request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

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

		newTask := models.Task{2, "Exercise"}
		jsonData, err := json.Marshal(newTask)
		testutils.AssertNoError(t, err)
		request := httptest.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

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
		request := httptest.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer([]byte(invalidJson)),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		testutils.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutils.AssertCalls(t, datastore.createTaskCalls, 0)
	})
}

var initialTasks = []models.Task{{1, "Pack clothes"}}

type mockDataStore struct {
	createTaskCalls int
	getTasksCalls   int
	tasks           []models.Task
}

func newMockDataStore() *mockDataStore {
	m := &mockDataStore{tasks: initialTasks}
	return m
}

func (m *mockDataStore) CreateTask(task models.Task) (models.Task, error) {
	m.createTaskCalls++
	return task, nil
}

func (m *mockDataStore) GetTasks() ([]models.Task, error) {
	m.getTasksCalls++
	return m.tasks, nil
}

func (m *mockDataStore) DeleteTaskById(id int) error {
	m.tasks = slices.DeleteFunc(m.tasks, func(task models.Task) bool {
		return task.Id == id
	})
	return nil
}
