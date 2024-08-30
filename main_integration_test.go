package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/api"
	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestServer(t *testing.T) {
	dbFile, cleanDatabase := testutils.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := data.NewFileSystemStore(dbFile)
	assert.HasNoError(t, err)
	server := api.NewServer(store)

	initialTasks := []models.Task{
		{1, "Buy groceries"},
		{2, "Pack clothes"},
	}

	for _, task := range initialTasks {
		_, err := sendPostTask(server, task)
		assert.HasNoError(t, err)
	}

	t.Run("responds with a 200 OK status on GET `/`", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusOK)
	})

	t.Run("responds with a 404 not found status on an invalid path", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/not-found", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns a slice of tasks with GET `/tasks`", func(t *testing.T) {
		response := sendGetTasks(server)
		assert.Status(t, response.Code, http.StatusOK)
		tasks := testutils.GetTasksFromResponse(t, response.Body)
		assert.Equals(t, tasks, initialTasks)
	})

	t.Run("returns the correct task with GET `/tasks/{id}`", func(t *testing.T) {
		wantedTask := initialTasks[1]
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/tasks/%d", wantedTask.Id),
			nil,
		)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusOK)
		got := testutils.GetTaskFromResponse(t, response.Body)
		assert.Equals(t, got, wantedTask)
	})

	t.Run("deletes the task with DELETE `/tasks/{id}`", func(t *testing.T) {
		newTask := models.Task{3, "Cook food"}
		postResponse, err := sendPostTask(server, newTask)
		assert.HasNoError(t, err)
		newTask = testutils.GetTaskFromResponse(t, postResponse.Body)

		deleteRequest := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/tasks/%d", newTask.Id),
			nil,
		)
		deleteResponse := httptest.NewRecorder()
		server.ServeHTTP(deleteResponse, deleteRequest)
		assert.Status(t, deleteResponse.Code, http.StatusNoContent)

		getResponse := sendGetTasks(server)
		tasks := testutils.GetTasksFromResponse(t, getResponse.Body)
		assert.DoesNotContain(t, tasks, newTask)
	})
}

func sendGetTasks(server *api.Server) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func sendPostTask(server *api.Server, task models.Task) (*httptest.ResponseRecorder, error) {
	jsonBody, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBuffer(jsonBody),
	)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response, nil
}
