package datastore_test

import (
	"slices"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/datastore"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

func TestFileSystemDataStore(t *testing.T) {
	wantedTasks := []models.Task{{Title: "Buy groceries"}}
	jsonTasks, err := utils.ConvertToJSON(wantedTasks)
	testutils.AssertNoError(t, err)
	database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
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
