package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestHandleGetTaskById(t *testing.T) {
	t.Run("returns the wanted task if it exists", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		wantedTask := data.Tasks[0]
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
			*testutils.GetTaskFromResponse(t, response.Body),
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
