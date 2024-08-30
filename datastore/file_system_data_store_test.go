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
	initialTasks := []models.Task{{Id: 1, Title: "Buy groceries"}}
	jsonTasks, err := utils.ConvertToJSON(initialTasks)
	testutils.AssertNoError(t, err)

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
	})

	t.Run("GetTasks returns the stored tasks", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
		tasks, err := store.GetTasks()
		testutils.AssertNoError(t, err)
		testutils.AssertEquals(t, tasks, initialTasks)
	})

	t.Run("CreateTask stores and returns the created task", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
		newTask := models.Task{Id: 2, Title: "Launder clothes"}
		store.CreateTask(newTask)
		tasks, err := store.GetTasks()
		testutils.AssertNoError(t, err)
		if !slices.Contains(tasks, newTask) {
			t.Errorf("missing task '%v' from tasks '%v'", newTask, tasks)
		}
	})

	t.Run("DeleteTaskById deletes the selected task", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := datastore.NewFileSystemDataStore(database)
		testutils.AssertNoError(t, err)
		store.DeleteTaskById(initialTasks[0].Id)
		tasks, err := store.GetTasks()
		testutils.AssertNoError(t, err)
		if slices.Contains(tasks, initialTasks[0]) {
			t.Errorf(
				"expected task '%+v' to be deleted but isn't",
				initialTasks[0],
			)
		}
	})

}
