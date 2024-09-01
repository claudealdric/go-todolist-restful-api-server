package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

func TestHandleRoot(t *testing.T) {
	t.Run("responds with 200 OK status", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusOK)
	})
}

func TestHandleDeleteTaskById(t *testing.T) {
	t.Run("deletes the task and responds with 204 No Content", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		tasks, err := data.GetTasks()
		assert.HasNoError(t, err)
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
		assert.Status(t, response.Code, http.StatusNoContent)

		tasks, err = data.GetTasks()
		assert.HasNoError(t, err)
		if len(tasks) != initialTasksCount-1 {
			t.Errorf(
				"expected a slice of length %d; received %+v",
				initialTasksCount-1,
				tasks,
			)
		}
	})

	t.Run("responds with a 400 Bad Request when sending a non-integer ID", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		request := httptest.NewRequest(
			http.MethodDelete,
			"/tasks/not-an-integer",
			nil,
		)
		response := httptest.NewRecorder()

		server.Handler.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusBadRequest)
	})

	t.Run("responds with 404 Not Found when the task does not exist", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		taskToDelete := models.Task{-1, "Does not exist"}
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/tasks/%d", taskToDelete.Id),
			nil,
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusNotFound)
	})

	t.Run("responds with 500 error when the store task deletion fails for an unknown reason", func(t *testing.T) {
		data := newMockStore(true)
		server := NewServer(data)

		taskToDelete := initialTasks[0]
		request := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/tasks/%d", taskToDelete.Id),
			nil,
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusInternalServerError)
	})
}

func TestHandleGetTaskById(t *testing.T) {
	t.Run("returns the wanted task if it exists", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		wantedTask := initialTasks[0]
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/tasks/%d", wantedTask.Id),
			nil,
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.ContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		assert.Status(t, response.Code, http.StatusOK)
		assert.Calls(t, data.getTaskByIdCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			wantedTask,
		)
	})

	t.Run("responds with a 400 Bad Request when given non-integer ID", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		invalidId := "not-an-integer"
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/tasks/%s", invalidId),
			nil,
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, data.getTaskByIdCalls, 0)
	})

	t.Run("responds with a 404 Not Found when the task cannot be found", func(t *testing.T) {
		data := newMockStore(true)
		server := NewServer(data)

		doesNotExistId := -1
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/tasks/%d", doesNotExistId),
			nil,
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusNotFound)
		assert.Calls(t, data.getTaskByIdCalls, 1)
	})

	t.Run("responds with a 500 error when the task retrieval fails for an unknown reason", func(t *testing.T) {
		data := newMockStore(true)
		server := NewServer(data)

		doesNotExistId := -2
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/tasks/%d", doesNotExistId),
			nil,
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusInternalServerError)
		assert.Calls(t, data.getTaskByIdCalls, 1)
	})
}

func TestHandleGetTasks(t *testing.T) {
	t.Run("returns the stored tasks", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.ContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		assert.Status(t, response.Code, http.StatusOK)
		assert.Calls(t, data.getTasksCalls, 1)
		assert.Equals(
			t,
			testutils.GetTasksFromResponse(t, response.Body),
			initialTasks,
		)
	})

	t.Run("responds with a 500 error when getting tasks from the store errors", func(t *testing.T) {
		data := newMockStore(true)
		server := NewServer(data)

		request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusInternalServerError)
		assert.Calls(t, data.getTasksCalls, 1)
	})
}

func TestHandlePostTasks(t *testing.T) {
	t.Run("creates and returns the task with a 201 Status Created", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		newTask := models.Task{2, "Exercise"}
		jsonData, err := json.Marshal(newTask)
		assert.HasNoError(t, err)
		request := httptest.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.ContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		assert.Status(t, response.Code, http.StatusCreated)
		assert.Calls(t, data.createTaskCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			newTask,
		)
	})

	t.Run("responds with a 400 Bad Request given an invalid body", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		invalidJson := `{`
		request := httptest.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer([]byte(invalidJson)),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, data.createTaskCalls, 0)
	})

	t.Run("responds with a 500 error when the store task creation fails", func(t *testing.T) {
		data := newMockStore(true)
		server := NewServer(data)

		newTask := models.Task{2, "Exercise"}
		jsonData, err := json.Marshal(newTask)
		assert.HasNoError(t, err)
		request := httptest.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusInternalServerError)
		assert.Calls(t, data.createTaskCalls, 1)
	})
}

func TestHandlePatchTasks(t *testing.T) {
	t.Run("returns the updated task and responds with a 200 OK status", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		task := initialTasks[0]

		newTitle := "Pack bags"
		dto := models.UpdateTaskDTO{Title: &newTitle}
		jsonData, err := json.Marshal(dto)
		assert.HasNoError(t, err)

		request := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/tasks/%d", task.Id),
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.ContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		assert.Status(t, response.Code, http.StatusOK)
		assert.Calls(t, data.updateTaskCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			models.Task{Id: task.Id, Title: newTitle},
		)
	})

	t.Run("responds with a 400 Bad Request with an invalid ID", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		invalidId := "not-an-integer"
		newTitle := "Pack bags"
		dto := models.UpdateTaskDTO{Title: &newTitle}
		jsonData, err := json.Marshal(dto)
		assert.HasNoError(t, err)

		request := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/tasks/%s", invalidId),
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, data.updateTaskCalls, 0)
	})

	t.Run("responds with a 404 Not Found when the task cannot be found", func(t *testing.T) {
		data := newMockStore(false)
		server := NewServer(data)

		doesNotExistId := -1
		newTitle := "Pack bags"
		dto := models.UpdateTaskDTO{Title: &newTitle}
		jsonData, err := json.Marshal(dto)
		assert.HasNoError(t, err)

		request := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/tasks/%d", doesNotExistId),
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusNotFound)
		assert.Calls(t, data.updateTaskCalls, 1)
	})
}

var initialTasks = []models.Task{{1, "Pack clothes"}}

type mockStore struct {
	createTaskCalls  int
	getTaskByIdCalls int
	getTasksCalls    int
	updateTaskCalls  int
	tasks            []models.Task
	shouldError      bool
}

func newMockStore(shouldError bool) *mockStore {
	m := &mockStore{tasks: initialTasks, shouldError: shouldError}
	return m
}

func (m *mockStore) CreateTask(task models.Task) (models.Task, error) {
	m.createTaskCalls++
	if m.shouldError {
		return models.Task{}, errors.New("forced error")
	}
	return task, nil
}

func (m *mockStore) GetTaskById(id int) (models.Task, error) {
	m.getTaskByIdCalls++
	var task models.Task
	if m.shouldError {
		if id == -1 {
			return task, data.ErrResourceNotFound
		} else {
			return task, errors.New("forced error")
		}
	}
	tasks, _ := m.GetTasks()
	task, _ = utils.SliceFind(tasks, func(t models.Task) bool {
		return t.Id == id
	})
	return task, nil
}

func (m *mockStore) GetTasks() ([]models.Task, error) {
	m.getTasksCalls++
	if m.shouldError {
		return nil, errors.New("forced error")
	}
	return m.tasks, nil
}

func (m *mockStore) DeleteTaskById(id int) error {
	if m.shouldError {
		return errors.New("forced error")
	}
	i := slices.IndexFunc(m.tasks, func(task models.Task) bool {
		return task.Id == id
	})
	if i == -1 {
		return data.ErrResourceNotFound
	}
	m.tasks = slices.DeleteFunc(m.tasks, func(task models.Task) bool {
		return task.Id == id
	})
	return nil
}

func (m *mockStore) UpdateTask(task models.Task) (models.Task, error) {
	m.updateTaskCalls++
	for i, t := range m.tasks {
		if t.Id == task.Id {
			m.tasks[i] = task
			return task, nil
		}
	}
	return models.Task{}, data.ErrResourceNotFound
}
