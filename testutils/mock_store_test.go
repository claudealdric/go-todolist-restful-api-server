package testutils

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestGetUserByEmail(t *testing.T) {
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

}
