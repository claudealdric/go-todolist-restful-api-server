package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestServer(t *testing.T) {
	server := NewServer()

	t.Run("get root", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Errorf("an error occurred during the request: %v", err)
		}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		testutils.AssertStatus(t, response.Code, http.StatusOK)
	})
}
