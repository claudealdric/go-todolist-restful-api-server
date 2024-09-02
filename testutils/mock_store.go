package testutils

import (
	"errors"
	"slices"

	"github.com/claudealdric/go-todolist-restful-api-server/data"
	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

var initialMockStoreTasks = []models.Task{{1, "Pack clothes"}}
var forcedError = errors.New("forced error")

type mockStore struct {
	CreateTaskCalls  int
	GetTaskByIdCalls int
	GetTasksCalls    int
	UpdateTaskCalls  int
	tasks            []models.Task
	shouldForceError bool
}

func NewMockStore(shouldError bool) *mockStore {
	m := &mockStore{tasks: initialMockStoreTasks, shouldForceError: shouldError}
	return m
}

func (m *mockStore) CreateTask(task models.Task) (models.Task, error) {
	m.CreateTaskCalls++
	if m.shouldForceError {
		return models.Task{}, forcedError
	}
	m.tasks = append(m.tasks, task)
	return task, nil
}

func (m *mockStore) GetInitialTasks() []models.Task {
	return initialMockStoreTasks
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
	return m.tasks, nil
}

func (m *mockStore) DeleteTaskById(id int) error {
	if m.shouldForceError {
		return forcedError
	}
	i := slices.IndexFunc(m.tasks, func(task models.Task) bool {
		return task.Id == id
	})
	if i == -1 {
		return data.ErrResourceNotFound
	}
	m.tasks = slices.DeleteFunc(m.tasks, func(task models.Task) bool {
		return task.Id == id
	})
	return nil
}

func (m *mockStore) UpdateTask(task models.Task) (models.Task, error) {
	m.UpdateTaskCalls++
	if m.shouldForceError {
		return models.Task{}, forcedError
	}
	for i, t := range m.tasks {
		if t.Id == task.Id {
			m.tasks[i] = task
			return task, nil
		}
	}
	return models.Task{}, data.ErrResourceNotFound
}
