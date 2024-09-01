package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestHandleRoot(t *testing.T) {
	t.Run("responds with 200 OK status", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusOK)
	})
}

func TestHandleDeleteTaskById(t *testing.T) {
	t.Run("deletes the task and responds with 204 No Content", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		data := testutils.NewMockStore(false)
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
		data := testutils.NewMockStore(false)
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
		data := testutils.NewMockStore(true)
		server := NewServer(data)

		taskToDelete := testutils.InitialMockStoreTasks[0]
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
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		wantedTask := testutils.InitialMockStoreTasks[0]
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
		assert.Calls(t, data.GetTaskByIdCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			wantedTask,
		)
	})

	t.Run("responds with a 400 Bad Request when given non-integer ID", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		assert.Calls(t, data.GetTaskByIdCalls, 0)
	})

	t.Run("responds with a 404 Not Found when the task cannot be found", func(t *testing.T) {
		data := testutils.NewMockStore(true)
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
		assert.Calls(t, data.GetTaskByIdCalls, 1)
	})

	t.Run("responds with a 500 error when the task retrieval fails for an unknown reason", func(t *testing.T) {
		data := testutils.NewMockStore(true)
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
		assert.Calls(t, data.GetTaskByIdCalls, 1)
	})
}

func TestHandleGetTasks(t *testing.T) {
	t.Run("returns the stored tasks", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		assert.Calls(t, data.GetTasksCalls, 1)
		assert.Equals(
			t,
			testutils.GetTasksFromResponse(t, response.Body),
			testutils.InitialMockStoreTasks,
		)
	})

	t.Run("responds with a 500 error when getting tasks from the store errors", func(t *testing.T) {
		data := testutils.NewMockStore(true)
		server := NewServer(data)

		request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusInternalServerError)
		assert.Calls(t, data.GetTasksCalls, 1)
	})
}

func TestHandlePostTasks(t *testing.T) {
	t.Run("creates and returns the task with a 201 Status Created", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		assert.Calls(t, data.CreateTaskCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			newTask,
		)
	})

	t.Run("responds with a 400 Bad Request given an invalid body", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		assert.Calls(t, data.CreateTaskCalls, 0)
	})

	t.Run("responds with a 500 error when the store task creation fails", func(t *testing.T) {
		data := testutils.NewMockStore(true)
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
		assert.Calls(t, data.CreateTaskCalls, 1)
	})
}

func TestHandlePatchTasks(t *testing.T) {
	t.Run("returns the updated task and responds with a 200 OK status", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		task := testutils.InitialMockStoreTasks[0]

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
		assert.Calls(t, data.UpdateTaskCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			models.Task{Id: task.Id, Title: newTitle},
		)
	})

	t.Run("responds with a 400 Bad Request with an invalid ID", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		assert.Calls(t, data.UpdateTaskCalls, 0)
	})

	t.Run("responds with a 404 Not Found when the task cannot be found", func(t *testing.T) {
		data := testutils.NewMockStore(false)
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
		assert.Calls(t, data.UpdateTaskCalls, 1)
	})

	t.Run("responds with a 400 Bad Request when the body is invalid", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		task := testutils.InitialMockStoreTasks[0]

		request := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/tasks/%d", task.Id),
			bytes.NewBuffer([]byte(`{`)),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, data.UpdateTaskCalls, 0)
	})

	t.Run("responds with a 500 error when an unknown store error occurs", func(t *testing.T) {
		data := testutils.NewMockStore(true)
		server := NewServer(data)

		task := testutils.InitialMockStoreTasks[0]

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

		assert.Status(t, response.Code, http.StatusInternalServerError)
		assert.Calls(t, data.UpdateTaskCalls, 1)
	})

	t.Run("providing the ID in the request body does not override the ID URL param", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)
		unmodifiedTask := models.Task{Id: 2, Title: "Exercise"}
		server.store.CreateTask(unmodifiedTask)

		taskToUpdate := testutils.InitialMockStoreTasks[0]

		newTitle := "Pack bags"
		jsonData, err := json.Marshal(models.Task{Id: 2, Title: newTitle})
		assert.HasNoError(t, err)

		request := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/tasks/%d", taskToUpdate.Id),
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
		assert.Calls(t, data.UpdateTaskCalls, 1)
		assert.Equals(
			t,
			testutils.GetTaskFromResponse(t, response.Body),
			models.Task{Id: taskToUpdate.Id, Title: newTitle},
		)

		gotUnmodifiedTask, _ := server.store.GetTaskById(unmodifiedTask.Id)
		assert.Equals(t, gotUnmodifiedTask, unmodifiedTask)
	})
}
