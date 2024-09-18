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
	"golang.org/x/crypto/bcrypt"
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
		userFromResponse := testutils.GetUserFromResponse(t, response.Body)
		gotUser := models.User{
			Id:    userFromResponse.Id,
			Name:  userFromResponse.Name,
			Email: userFromResponse.Email,
		}
		newUser := models.NewUser(1, dto.Name, dto.Email, dto.Password)
		wantedUser := models.User{
			Id:    newUser.Id,
			Name:  newUser.Name,
			Email: newUser.Email,
		}

		assert.HasNoError(t, bcrypt.CompareHashAndPassword(
			[]byte(userFromResponse.Password),
			[]byte(newUser.Password),
		))
		assert.ContentType(
			t,
			testutils.GetContentTypeFromResponse(response),
			jsonContentType,
		)
		assert.Status(t, response.Code, http.StatusCreated)
		assert.Calls(t, data.CreateUserCalls, 1)
		assert.Equals(t, gotUser, wantedUser)
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
