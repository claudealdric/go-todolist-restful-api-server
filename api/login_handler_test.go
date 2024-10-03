package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestHandleLogin(t *testing.T) {
	t.Run("returns a 400 Bad Request status when given an invalid body", func(t *testing.T) {
		store := testutils.NewMockStore(true)
		server := NewServer(store)

		invalidJson := `{`
		request := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer([]byte(invalidJson)),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, store.ValidateUserCredentialsCalls, 0)
	})

	t.Run("returns a 400 Bad Request status when given an empty body", func(t *testing.T) {
		store := testutils.NewMockStore(true)
		server := NewServer(store)

		request := httptest.NewRequest(http.MethodPost, "/login", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, store.ValidateUserCredentialsCalls, 0)
	})

	t.Run("returns a 401 Unauthorized status when given incorrect credentials", func(t *testing.T) {
		store := testutils.NewMockStore(true)
		server := NewServer(store)

		credentials := LoginCredentials{
			Email:    "does-not-exist@mail.com",
			Password: "password",
		}
		jsonData, err := json.Marshal(credentials)
		assert.HasNoError(t, err)

		request := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusUnauthorized)
		assert.Calls(t, store.ValidateUserCredentialsCalls, 1)
	})

	t.Run("returns a 200 OK status when given the correct credentials", func(t *testing.T) {
		store := testutils.NewMockStore(false)
		server := NewServer(store)

		credentials := LoginCredentials{
			Email:    "the1@email.com",
			Password: "password",
		}
		jsonData, err := json.Marshal(credentials)
		assert.HasNoError(t, err)

		request := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusOK)
		assert.Calls(t, store.ValidateUserCredentialsCalls, 1)
	})
}
