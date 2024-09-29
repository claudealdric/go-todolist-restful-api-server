package data_test

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
	"golang.org/x/crypto/bcrypt"
)

func TestFileSystemStoreTasks(t *testing.T) {
	initialTasks := []models.Task{*models.NewTask(1, "Buy groceries")}
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
		assert.Equals(t, *got, wantedTask)
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

		newTask, err := store.CreateTask(&models.CreateTaskDTO{"Launder clothes"})
		assert.HasNoError(t, err)

		tasks, err := store.GetTasks()
		assert.HasNoError(t, err)

		assert.Contains(t, tasks, *newTask)
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
			&models.Task{
				Id:    task.Id,
				Title: newTitle,
			},
		)
		assert.HasNoError(t, err)
		wantedTask := models.Task{Id: task.Id, Title: newTitle}
		assert.Equals(t, *updatedTask, wantedTask)

		retrievedTask, err := store.GetTaskById(task.Id)
		assert.HasNoError(t, err)
		assert.Equals(t, *retrievedTask, wantedTask)
	})

	t.Run("UpdateTaskById returns an error with an invalid ID", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonTasks))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)
		assert.HasNoError(t, err)

		newTitle := "Buy food"
		_, err = store.UpdateTask(
			&models.Task{
				Id:    -1,
				Title: newTitle,
			},
		)
		assert.HasError(t, err)
	})
}

func TestFileSystemStoreUsers(t *testing.T) {
	initialUsers := []models.User{
		models.User{
			Id:       1,
			Name:     "Claude Aldric",
			Email:    "claude.aldric@email.com",
			Password: "password",
		},
	}
	jsonUsers, err := utils.ConvertToJSON(initialUsers)
	assert.HasNoError(t, err)

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)
	})

	t.Run("CreateUser stores and returns the created user", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, "")
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		dto := models.CreateUserDTO{
			Name:     "John Doe",
			Email:    "john.doe@email.com",
			Password: "password",
		}
		newUser, err := store.CreateUser(&dto)
		gotUser := models.User{
			Id:    newUser.Id,
			Name:  newUser.Name,
			Email: newUser.Email,
		}
		wantedUser := models.User{
			Id:    1,
			Name:  dto.Name,
			Email: dto.Email,
		}
		assert.HasNoError(t, err)
		assert.HasNoError(t, bcrypt.CompareHashAndPassword(
			[]byte(newUser.Password),
			[]byte(dto.Password),
		))
		assert.Equals(t, gotUser, wantedUser)

		users, err := store.GetUsers()
		assert.HasNoError(t, err)
		assert.Contains(t, users, *newUser)
	})

	t.Run("GetUserByEmail returns the correct user if it exists", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonUsers))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)
		assert.HasNoError(t, err)

		wantedUser := initialUsers[0]
		got, err := store.GetUserByEmail(wantedUser.Email)

		assert.HasNoError(t, err)
		assert.Equals(t, *got, wantedUser)
	})

	t.Run("GetUserByEmail returns an `ErrNotFound` error if user does not exist", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonUsers))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)
		assert.HasNoError(t, err)

		doesNotExistEmail := "does.not@exist.com"
		_, err = store.GetUserByEmail(doesNotExistEmail)

		assert.ErrorContains(t, err, data.ErrResourceNotFound)
	})

	t.Run("GetUsers returns the stored users", func(t *testing.T) {
		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonUsers))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)

		assert.HasNoError(t, err)

		users, err := store.GetUsers()

		assert.HasNoError(t, err)
		assert.Equals(t, users, initialUsers)
	})

	t.Run("ValidateUserCredentials", func(t *testing.T) {
		createUserDTO := models.NewCreateUserDTO(
			"Claude Aldric",
			"cvaldric@gmail.com",
			"Caput Draconis",
		)
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(createUserDTO.Password),
			bcrypt.DefaultCost,
		)
		assert.HasNoError(t, err)

		initialUsers := []models.User{
			models.User{
				Id:       1,
				Name:     "Claude Aldric",
				Email:    "cvaldric@gmail.com",
				Password: string(hashedPassword),
			},
		}
		jsonUsers, err := utils.ConvertToJSON(initialUsers)
		assert.HasNoError(t, err)

		database, cleanDatabase := testutils.CreateTempFile(t, string(jsonUsers))
		defer cleanDatabase()

		store, err := data.NewFileSystemStore(database)
		assert.HasNoError(t, err)

		tests := []struct {
			name     string
			email    string
			password string
			want     bool
		}{
			{
				name:     "when provided valid credentials",
				email:    "cvaldric@gmail.com",
				password: "Caput Draconis",
				want:     true,
			},
			{
				name:     "when provided an invalid password",
				email:    "cvaldric@gmail.com",
				password: "password",
				want:     false,
			},
			{
				name:     "when provided an invalid email",
				email:    "doesnotexist@email.com",
				password: "Caput Draconis",
				want:     false,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				assert.Equals(
					t,
					store.ValidateUserCredentials(test.email, test.password),
					test.want,
				)
			})
		}
	})
}
