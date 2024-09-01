package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

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
