package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

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
			data.Tasks,
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
