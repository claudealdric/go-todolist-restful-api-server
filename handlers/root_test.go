package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/handlers"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestHandleRoot(t *testing.T) {
	t.Run("responds with 200 OK status on the root path", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HandleRoot)
		handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("responds with a 404 not found with an invalid path", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/not-found", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HandleRoot)
		handler.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}
