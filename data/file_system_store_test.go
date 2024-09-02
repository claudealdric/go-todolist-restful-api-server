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

	assert.HasNoError(t, err)

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)
	})

	t.Run("GetTasks returns the stored tasks", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		tasks, err := store.GetTasks()

		assert.HasNoError(t, err)
		assert.Equals(t, tasks, initialTasks)
	})

	t.Run("GetTaskById returns the correct task if it exists", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		wantedTask := initialTasks[0]
		got, err := store.GetTaskById(wantedTask.Id)

		assert.HasNoError(t, err)
		assert.Equals(t, got, wantedTask)
	})

	t.Run("GetTaskById returns an `ErrNotFound` error if task does not exist", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		doesNotExistId := -1
		_, err = store.GetTaskById(doesNotExistId)

		assert.ErrorContains(t, err, data.ErrResourceNotFound)
	})

	t.Run("CreateTask stores and returns the created task", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		newTask, err := store.CreateTask(models.CreateTaskDTO{"Launder clothes"})
		assert.HasNoError(t, err)

		tasks, err := store.GetTasks()
		assert.HasNoError(t, err)

		assert.Contains(t, tasks, newTask)
	})

	t.Run("DeleteTaskById deletes the selected task", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		taskToDelete := initialTasks[0]
		store.DeleteTaskById(taskToDelete.Id)
		tasks, err := store.GetTasks()

		assert.HasNoError(t, err)
		assert.DoesNotContain(t, tasks, taskToDelete)
	})

	t.Run("DeleteTaskById returns with a `ErrResourceNotFound` error if task does not exist", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		doesNotExistId := -1
		err = store.DeleteTaskById(doesNotExistId)

		assert.ErrorContains(t, err, data.ErrResourceNotFound)
	})

	t.Run("UpdateTaskById updates and returns the task if it exists", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)
		assert.HasNoError(t, err)

		task := initialTasks[0]
		newTitle := "Buy food"
		updatedTask, err := store.UpdateTask(
			models.Task{
				Id:    task.Id,
				Title: newTitle,
			},
		)
		assert.HasNoError(t, err)
		wantedTask := models.Task{Id: task.Id, Title: newTitle}
		assert.Equals(t, updatedTask, wantedTask)

		retrievedTask, err := store.GetTaskById(task.Id)
		assert.HasNoError(t, err)
		assert.Equals(t, retrievedTask, wantedTask)
	})

	t.Run("UpdateTaskById returns an error with an invalid ID", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)
		assert.HasNoError(t, err)

		newTitle := "Buy food"
		_, err = store.UpdateTask(
			models.Task{
				Id:    -1,
				Title: newTitle,
			},
		)
		assert.HasError(t, err)
	})
}
