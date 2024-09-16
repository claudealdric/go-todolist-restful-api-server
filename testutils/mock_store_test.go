package testutils

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestGetUserByEmail(t *testing.T) {
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
}
