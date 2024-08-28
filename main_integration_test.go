package main

import (
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
	// TODO: start with an empty database
	dbFile, cleanDatabase := testutils.CreateTempFile(
		t,
		`[{"title": "Buy groceries"}]`,
	)
	defer cleanDatabase()
	store, err := datastore.NewFileSystemDataStore(dbFile)
	testutils.AssertNoError(t, err)
	server := server.NewServer(store)

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
		want := []models.Task{{Title: "Buy groceries"}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
