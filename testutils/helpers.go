package testutils

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

func AssertCalls(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect number of calls; got %d, want %d", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}

}

func AssertEquals[T any](t testing.TB, got, want T) {
	t.Helper()
	switch v := any(got).(type) {
	case string, int, int64, float64, bool:
		if v != any(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	default:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func CreateTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tempFile.Write([]byte(initialData))

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile
}

func GetTaskFromResponse(t *testing.T, body io.Reader) (tasks models.Task) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&tasks)

	if err != nil {
		t.Fatalf(
			"unable to parse response from server %q into Task: %v",
			body,
			err,
		)
	}

	return tasks
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
