package testutils

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestGetUserByEmail(t *testing.T) {
	t.Run("CreateUser increments the internal counter and returns the new user", func(t *testing.T) {
		mockStore := NewMockStore(false)
		dto := models.CreateUserDTO{
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		}
		wantedUser := models.User{
			Id:       1,
			Name:     dto.Name,
			Email:    dto.Email,
			Password: dto.Password,
		}
		gotUser, err := mockStore.CreateUser(dto)

		assert.HasNoError(t, err)
		assert.Equals(t, mockStore.CreateUserCalls, 1)
		assert.Equals(t, *gotUser, wantedUser)
	})

	t.Run("forcing CreateUser to fail returns the forced error", func(t *testing.T) {
		mockStore := NewMockStore(true)
		dto := models.CreateUserDTO{
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		}
		gotUser, err := mockStore.CreateUser(dto)

		assert.ErrorContains(t, err, forcedError)
		assert.Equals(t, mockStore.CreateUserCalls, 1)
		assert.Equals(t, gotUser, nil)
	})

	t.Run("GetUserByEmail increments the internal counter and returns the wanted user", func(t *testing.T) {
		mockStore := NewMockStore(false)
		wantedUser := models.User{
			Id:       1,
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		}
		mockStore.Users = []models.User{wantedUser}
		gotUser, err := mockStore.GetUserByEmail(wantedUser.Email)

		assert.HasNoError(t, err)
		assert.Equals(t, mockStore.GetUserByEmailCalls, 1)
		assert.Equals(t, gotUser, wantedUser)
	})

	t.Run("forcing GetUserByEmail to fail returns the forced error", func(t *testing.T) {
		mockStore := NewMockStore(true)
		wantedUser := models.User{
			Id:       1,
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		}
		mockStore.Users = []models.User{wantedUser}
		gotUser, err := mockStore.GetUserByEmail(wantedUser.Email)

		assert.ErrorContains(t, err, forcedError)
		assert.Equals(t, gotUser, models.User{})
		assert.Equals(t, mockStore.GetUserByEmailCalls, 1)
	})

	t.Run("GetUsers returns all users", func(t *testing.T) {
		mockStore := NewMockStore(false)
		initialUsers := []models.User{
			models.User{
				Id:       1,
				Name:     "Claude Aldric",
				Email:    "claude.aldric@email.com",
				Password: "password",
			},
		}
		mockStore.Users = initialUsers
		users, err := mockStore.GetUsers()

		assert.HasNoError(t, err)
		assert.Equals(t, mockStore.GetUsersCalls, 1)
		assert.Equals(t, users, initialUsers)
	})

	t.Run("forcing GetUsers to fail returns the forced error", func(t *testing.T) {
		mockStore := NewMockStore(true)
		initialUsers := []models.User{
			models.User{
				Id:       1,
				Name:     "Claude Aldric",
				Email:    "claude.aldric@email.com",
				Password: "password",
			},
		}
		mockStore.Users = initialUsers
		users, err := mockStore.GetUsers()

		assert.ErrorContains(t, err, forcedError)
		assert.Equals(t, users, nil)
		assert.Equals(t, mockStore.GetUsersCalls, 1)
	})
}
