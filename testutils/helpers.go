package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

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

func GetContentTypeFromResponse(response *httptest.ResponseRecorder) string {
	return response.Result().Header.Get("content-type")
}

func GetTaskFromResponse(t *testing.T, body io.Reader) *models.Task {
	t.Helper()
	var tasks models.Task
	err := json.NewDecoder(body).Decode(&tasks)

	if err != nil {
		fmt.Printf(
			"%s: unable to parse response from server %q into Task: %v\n",
			utils.GetCurrentFunctionName(),
			body,
			err,
		)
		return nil
	}

	return &tasks
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

func GetUserFromResponse(t *testing.T, body io.Reader) (user models.User) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&user)

	if err != nil {
		t.Fatalf(
			"unable to parse response from server %q into User: %v",
			body,
			err,
		)
	}

	return user
}
