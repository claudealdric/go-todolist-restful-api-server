package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/server"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestServer(t *testing.T) {
	dbFile, cleanDatabase := testutils.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := datastore.NewFileSystemDataStore(dbFile)
	testutils.AssertNoError(t, err)
	server := server.NewServer(store)

	initialTasks := []models.Task{
		{1, "Buy groceries"},
		{2, "Pack clothes"},
	}

	for _, task := range initialTasks {
		_, err := sendPostTask(server, task)
		testutils.AssertNoError(t, err)
	}

	t.Run("responds with a 200 OK status on GET `/`", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("responds with a 404 not found status on an invalid path", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/not-found", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns a slice of tasks with GET `/tasks`", func(t *testing.T) {
		response := sendGetTasks(server)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
		tasks := testutils.GetTasksFromResponse(t, response.Body)
		testutils.AssertEquals(t, tasks, initialTasks)
	})

	t.Run("deletes the task with DELETE `/tasks/{id}`", func(t *testing.T) {
		newTask := models.Task{3, "Cook food"}
		postResponse, err := sendPostTask(server, newTask)
		testutils.AssertNoError(t, err)
		newTask = testutils.GetTaskFromResponse(t, postResponse.Body)

		deleteRequest := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/tasks/%d", newTask.Id),
			nil,
		)
		deleteResponse := httptest.NewRecorder()
		server.ServeHTTP(deleteResponse, deleteRequest)
		testutils.AssertStatus(t, deleteResponse.Code, http.StatusNoContent)

		getResponse := sendGetTasks(server)
		tasks := testutils.GetTasksFromResponse(t, getResponse.Body)
		testutils.AssertDoesNotContain(t, tasks, newTask)
	})
	// TODO: return a 404 Not Found if the resource cannot be found
}

func sendGetTasks(server *server.Server) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func sendPostTask(server *server.Server, task models.Task) (*httptest.ResponseRecorder, error) {
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
