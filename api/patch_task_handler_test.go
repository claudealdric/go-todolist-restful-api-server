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

func TestHandlePatchTask(t *testing.T) {
	t.Run("returns the updated task and responds with a 200 OK status", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		task := data.GetInitialTasks()[0]

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

		task := data.GetInitialTasks()[0]

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

		task := data.GetInitialTasks()[0]

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
		unmodifiedTask := models.Task{Id: 2, Title: "Exercise"}
		data.Tasks = append(data.Tasks, unmodifiedTask)
		server := NewServer(data)

		taskToUpdate := data.GetInitialTasks()[0]

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
