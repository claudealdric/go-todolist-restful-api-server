package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/api"
	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils"
	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestServerWithSqliteStore(t *testing.T) {
	dbFile := "./tmp/data.db"
	db, err := sql.Open("sqlite3", dbFile)
	assert.HasNoError(t, err)
	defer db.Close()
	data.InitDb(db)
	defer cleanSqliteDatabase(dbFile)
	store := data.NewSqliteStore(db)
	server := api.NewServer(store)

	t.Run("tasks", func(t *testing.T) {
		createTaskDTO := models.NewCreateTaskDTO("Write integration tests")
		createTaskResponse, err := sendPostTask(server, createTaskDTO)
		assert.HasNoError(t, err)
		createdTask := testutils.GetTaskFromResponse(t, createTaskResponse.Body)
		wantedTask := models.NewTask(createdTask.Id, createTaskDTO.Title)
		assert.Equals(t, createdTask, wantedTask)

		getTaskByIdResponse := sendGetTaskById(server, createdTask.Id)
		task := testutils.GetTaskFromResponse(t, getTaskByIdResponse.Body)
		assert.Equals(t, task, wantedTask)

		getTasksResponse := sendGetTasks(server)
		tasks := testutils.GetTasksFromResponse(t, getTasksResponse.Body)
		assert.Contains(t, tasks, *wantedTask)

		updatedTitle := "Profit"
		updateTaskDTO := models.UpdateTaskDTO{Title: &updatedTitle}
		patchTaskResponse, err := sendPatchTask(server, updateTaskDTO, createdTask.Id)
		assert.HasNoError(t, err)
		task = testutils.GetTaskFromResponse(t, patchTaskResponse.Body)
		wantedTask = models.NewTask(createdTask.Id, updatedTitle)
		assert.Equals(t, task, wantedTask)

		sendDeleteTask(server, createdTask.Id)
		unwantedTask := wantedTask

		getTaskByIdResponse = sendGetTaskById(server, createdTask.Id)
		task = testutils.GetTaskFromResponse(t, getTaskByIdResponse.Body)
		assert.Equals(t, task, nil)

		getTasksResponse = sendGetTasks(server)
		tasks = testutils.GetTasksFromResponse(t, getTasksResponse.Body)
		fmt.Println("tasks", tasks)
		assert.DoesNotContain(t, tasks, *unwantedTask)
	})

	t.Run("users", func(t *testing.T) {
		createUserDTO := models.NewCreateUserDTO(
			"Sherlock",
			"sherlock@email.com",
			"sherlocked",
		)
		postUserResponse, err := sendPostUser(server, createUserDTO)
		assert.HasNoError(t, err)
		createdUser := testutils.GetUserFromResponse(t, postUserResponse.Body)
		wantedUser := models.NewUser(
			createdUser.Id,
			createUserDTO.Name,
			createdUser.Email,
			createUserDTO.Password,
		)
		assert.Equals(t, createdUser, *wantedUser)
	})
}

func TestServerWithFileSystemStore(t *testing.T) {
	dbFile, cleanDatabase := testutils.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := data.NewFileSystemStore(dbFile)
	assert.HasNoError(t, err)
	server := api.NewServer(store)

	initialTasks := []models.Task{
		*models.NewTask(1, "Buy groceries"),
		*models.NewTask(2, "Pack clothes"),
	}

	for _, task := range initialTasks {
		dto := models.NewCreateTaskDTO(task.Title)
		_, err := sendPostTask(server, dto)
		assert.HasNoError(t, err)
	}

	t.Run("responds with a 200 OK status on GET `/`", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusOK)
	})

	t.Run("responds with a 404 not found status on an invalid path", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/not-found", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns a slice of tasks with GET `/tasks`", func(t *testing.T) {
		response := sendGetTasks(server)
		assert.Status(t, response.Code, http.StatusOK)
		tasks := testutils.GetTasksFromResponse(t, response.Body)
		assert.Equals(t, tasks, initialTasks)
	})

	t.Run("returns the correct task with GET `/tasks/{id}`", func(t *testing.T) {
		wantedTask := initialTasks[1]
		request := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/tasks/%d", wantedTask.Id),
			nil,
		)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Status(t, response.Code, http.StatusOK)
		got := testutils.GetTaskFromResponse(t, response.Body)
		assert.Equals(t, *got, wantedTask)
	})

	t.Run("deletes the task with DELETE `/tasks/{id}`", func(t *testing.T) {
		newTaskDto := models.NewCreateTaskDTO("Cook food")
		postResponse, err := sendPostTask(server, newTaskDto)
		assert.HasNoError(t, err)
		newTask := testutils.GetTaskFromResponse(t, postResponse.Body)

		deleteRequest := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/tasks/%d", newTask.Id),
			nil,
		)
		deleteResponse := httptest.NewRecorder()
		server.ServeHTTP(deleteResponse, deleteRequest)
		assert.Status(t, deleteResponse.Code, http.StatusNoContent)

		getResponse := sendGetTasks(server)
		tasks := testutils.GetTasksFromResponse(t, getResponse.Body)
		assert.DoesNotContain(t, tasks, *newTask)
	})

	t.Run("updates the task with PATCH `/tasks/{id}`", func(t *testing.T) {
		newTaskDto := models.NewCreateTaskDTO("Walk the dog")
		postResponse, err := sendPostTask(server, newTaskDto)
		assert.HasNoError(t, err)

		taskId := testutils.GetTaskFromResponse(t, postResponse.Body).Id

		newTitle := "Walk the cat"
		updateTaskDTO := models.UpdateTaskDTO{Title: &newTitle}
		patchResponse, err := sendPatchTask(server, updateTaskDTO, taskId)
		assert.HasNoError(t, err)

		wantedTask := models.NewTask(taskId, *updateTaskDTO.Title)

		updatedTask := testutils.GetTaskFromResponse(t, patchResponse.Body)
		assert.Status(t, patchResponse.Code, http.StatusOK)
		assert.Equals(t, updatedTask, wantedTask)

		getResponse := sendGetTaskById(server, taskId)
		task := testutils.GetTaskFromResponse(t, getResponse.Body)
		assert.Equals(t, task, wantedTask)
	})
}

func sendGetTaskById(server *api.Server, id int) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/tasks/%d", id),
		nil,
	)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func sendDeleteTask(
	server *api.Server,
	taskId int,
) {
	request := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/tasks/%d", taskId),
		nil,
	)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
}

func sendGetTasks(server *api.Server) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

func sendPatchTask(
	server *api.Server,
	DTO models.UpdateTaskDTO,
	taskId int,
) (
	*httptest.ResponseRecorder,
	error,
) {
	jsonBody, err := json.Marshal(DTO)
	if err != nil {
		return nil, err
	}
	request := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/tasks/%d", taskId),
		bytes.NewBuffer(jsonBody),
	)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response, nil
}

func sendPostTask(server *api.Server, dto *models.CreateTaskDTO) (
	*httptest.ResponseRecorder,
	error,
) {
	jsonBody, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBuffer(jsonBody),
	)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response, nil
}

func sendPostUser(server *api.Server, dto *models.CreateUserDTO) (
	*httptest.ResponseRecorder,
	error,
) {
	jsonBody, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewBuffer(jsonBody),
	)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response, nil
}

func cleanSqliteDatabase(path string) {
	os.Remove(path)
}
