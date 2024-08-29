package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	wantedTasks := []models.Task{
		{Title: "Buy groceries"},
		{Title: "Pack clothes"},
	}

	for _, task := range wantedTasks {
		jsonBody, err := json.Marshal(task)
		testutils.AssertNoError(t, err)
		request, err := http.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewBuffer(jsonBody),
		)
		testutils.AssertNoError(t, err)
		server.ServeHTTP(httptest.NewRecorder(), request)
	}

	t.Run("responds with a 200 OK status on the root path", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		testutils.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("responds with a 404 not found status on the root path", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/not-found", nil)
		testutils.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns a slice of tasks with GET on `/tasks`", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/tasks", nil)
		testutils.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)

		got := testutils.GetTasksFromResponse(t, response.Body)

		if !reflect.DeepEqual(got, wantedTasks) {
			t.Errorf("got %+v, want %+v", got, wantedTasks)
		}
	})
}
