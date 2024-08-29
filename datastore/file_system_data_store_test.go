package datastore_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
)

func TestFileSystemDataStore(t *testing.T) {
	wantedTasks := []models.Task{{Title: "Buy groceries"}}
	jsonTasks, err := ConvertToJSON(wantedTasks)
	testutils.AssertNoError(t, err)
	database, cleanDatabase := testutils.CreateTempFile(t, jsonTasks)
	defer cleanDatabase()

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
	})

	t.Run("GetTasks returns the stored tasks", func(t *testing.T) {
		store, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
		testutils.AssertEquals(t, store.GetTasks(), wantedTasks)
	})

	t.Run("CreateTask stores and returns the created task", func(t *testing.T) {
		store, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
		newTask := models.Task{Title: "Launder clothes"}
		store.CreateTask(newTask)
		tasks := store.GetTasks()
		if !slices.Contains(tasks, newTask) {
			t.Errorf("missing task '%v' from tasks '%v'", newTask, tasks)
		}
	})

}

func ConvertToJSON(tasks []models.Task) (string, error) {
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
