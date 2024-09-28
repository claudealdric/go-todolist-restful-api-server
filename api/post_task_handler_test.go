package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestHandlePostTask(t *testing.T) {
	t.Run("creates and returns the task with a 201 Status Created", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		newTask := models.NewTask(2, "Exercise")
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

		newTask := models.NewTask(2, "Exercise")
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
