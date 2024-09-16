package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestHandlePostUser(t *testing.T) {
	t.Run("creates and returns the user with a 201 Status Created", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		dto := models.CreateUserDTO{
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		}
		jsonData, err := json.Marshal(dto)
		assert.HasNoError(t, err)
		request := httptest.NewRequest(
			http.MethodPost,
			"/users",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		newUser := models.NewUser(1, dto.Name, dto.Email, dto.Password)

		assert.ContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		assert.Status(t, response.Code, http.StatusCreated)
		assert.Calls(t, data.CreateUserCalls, 1)
		assert.Equals(
			t,
			testutils.GetUserFromResponse(t, response.Body),
			*newUser,
		)
	})

	t.Run("responds with a 400 Bad Request given an invalid body", func(t *testing.T) {
		data := testutils.NewMockStore(false)
		server := NewServer(data)

		invalidJson := `{`
		request := httptest.NewRequest(
			http.MethodPost,
			"/users",
			bytes.NewBuffer([]byte(invalidJson)),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusBadRequest)
		assert.Calls(t, data.CreateUserCalls, 0)
	})

	t.Run("responds with a 500 error when the store user creation fails", func(t *testing.T) {
		data := testutils.NewMockStore(true)
		server := NewServer(data)

		dto := models.CreateUserDTO{
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		}
		jsonData, err := json.Marshal(dto)
		assert.HasNoError(t, err)
		request := httptest.NewRequest(
			http.MethodPost,
			"/users",
			bytes.NewBuffer(jsonData),
		)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		assert.Status(t, response.Code, http.StatusInternalServerError)
		assert.Calls(t, data.CreateUserCalls, 1)
	})
}
