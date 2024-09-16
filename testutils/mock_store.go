package testutils

import (
	"errors"
	"slices"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

var initialMockStoreTasks = []models.Task{*models.NewTask(1, "Pack clothes")}
var forcedError = errors.New("forced error")

type mockStore struct {
	CreateTaskCalls     int
	CreateUserCalls     int
	GetTaskByIdCalls    int
	GetTasksCalls       int
	GetUserByEmailCalls int
	GetUsersCalls       int
	Tasks               []models.Task
	UpdateTaskCalls     int
	Users               []models.User
	lastTaskId          int
	lastUserId          int
	shouldForceError    bool
}

func NewMockStore(shouldError bool) *mockStore {
	m := &mockStore{
		Tasks:            initialMockStoreTasks,
		shouldForceError: shouldError,
		lastTaskId:       1,
	}
	return m
}

func (m *mockStore) CreateTask(dto models.CreateTaskDTO) (models.Task, error) {
	m.CreateTaskCalls++
	if m.shouldForceError {
		return models.Task{}, forcedError
	}
	task := models.Task{Id: m.getNewTaskId(), Title: dto.Title}
	m.Tasks = append(m.Tasks, task)
	return task, nil
}

func (m *mockStore) GetTaskById(id int) (models.Task, error) {
	m.GetTaskByIdCalls++
	var task models.Task
	if m.shouldForceError {
		if id == -1 {
			return task, data.ErrResourceNotFound
		} else {
			return task, forcedError
		}
	}
	tasks, _ := m.GetTasks()
	task, _ = utils.SliceFind(tasks, func(t models.Task) bool {
		return t.Id == id
	})
	return task, nil
}

func (m *mockStore) GetTasks() ([]models.Task, error) {
	m.GetTasksCalls++
	if m.shouldForceError {
		return nil, forcedError
	}
	return m.Tasks, nil
}

func (m *mockStore) DeleteTaskById(id int) error {
	if m.shouldForceError {
		return forcedError
	}
	i := slices.IndexFunc(m.Tasks, func(task models.Task) bool {
		return task.Id == id
	})
	if i == -1 {
		return data.ErrResourceNotFound
	}
	m.Tasks = slices.DeleteFunc(m.Tasks, func(task models.Task) bool {
		return task.Id == id
	})
	return nil
}

func (m *mockStore) UpdateTask(task models.Task) (models.Task, error) {
	m.UpdateTaskCalls++
	if m.shouldForceError {
		return models.Task{}, forcedError
	}
	for i, t := range m.Tasks {
		if t.Id == task.Id {
			m.Tasks[i] = task
			return task, nil
		}
	}
	return models.Task{}, data.ErrResourceNotFound
}

func (m *mockStore) CreateUser(dto models.CreateUserDTO) (models.User, error) {
	m.CreateUserCalls++
	if m.shouldForceError {
		return models.User{}, forcedError
	}
	user := models.User{
		Id:       m.getNewUserId(),
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
	m.Users = append(m.Users, user)
	return user, nil
}

func (m *mockStore) GetUserByEmail(email string) (models.User, error) {
	m.GetUserByEmailCalls++
	if m.shouldForceError {
		return models.User{}, forcedError
	}
	user, _ := utils.SliceFind(m.Users, func(u models.User) bool {
		return u.Email == email
	})
	return user, nil
}

func (m *mockStore) GetUsers() ([]models.User, error) {
	m.GetUsersCalls++
	if m.shouldForceError {
		return nil, forcedError
	}
	return m.Users, nil
}

func (m *mockStore) getNewTaskId() int {
	newTaskId := m.lastTaskId + 1
	m.lastTaskId++
	return newTaskId
}

func (m *mockStore) getNewUserId() int {
	newUserId := m.lastUserId + 1
	m.lastUserId++
	return newUserId
}
