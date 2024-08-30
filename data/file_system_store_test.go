package data_test

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

func TestFileSystemStore(t *testing.T) {
	initialTasks := []models.Task{{1, "Buy groceries"}}
	jsonTasks, err := utils.ConvertToJSON(initialTasks)

	assert.NoError(t, err)

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := data.NewFileSystemStore(database)

		assert.NoError(t, err)
	})

	t.Run("GetTasks returns the stored tasks", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.NoError(t, err)

		tasks, err := store.GetTasks()

		assert.NoError(t, err)
		assert.Equals(t, tasks, initialTasks)
	})

	t.Run("GetTaskById returns the correct task if it exists", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.NoError(t, err)

		wantedTask := initialTasks[0]
		got, err := store.GetTaskById(wantedTask.Id)

		assert.NoError(t, err)
		assert.Equals(t, got, wantedTask)
	})

	t.Run("GetTaskById returns an error if it does not exist", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.NoError(t, err)

		invalidId := -1
		_, err = store.GetTaskById(invalidId)

		assert.HasError(t, err)
	})

	t.Run("CreateTask stores and returns the created task", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.NoError(t, err)

		newTask := models.Task{2, "Launder clothes"}
		store.CreateTask(newTask)
		tasks, err := store.GetTasks()

		assert.NoError(t, err)
		assert.Contains(t, tasks, newTask)
	})

	t.Run("DeleteTaskById deletes the selected task", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.NoError(t, err)

		taskToDelete := initialTasks[0]
		store.DeleteTaskById(taskToDelete.Id)
		tasks, err := store.GetTasks()

		assert.NoError(t, err)
		assert.DoesNotContain(t, tasks, taskToDelete)
	})

}
