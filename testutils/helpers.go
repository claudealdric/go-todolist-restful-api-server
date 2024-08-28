package testutils

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func GetTasksFromResponse(t *testing.T, body io.Reader) (tasks []models.Task) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&tasks)

	if err != nil {
		t.Fatalf(
			"unable to parse response from server %q into slice of Task: %v",
			body,
			err,
		)
	}

	return tasks
}
