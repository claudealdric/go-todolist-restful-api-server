package datastore

import "github.com/claudealdric/go-todolist-restful-api-server/models"

type DataStore interface {
	GetTasks() []models.Task
	CreateTask(task models.Task) (models.Task, error)
	DeleteTaskById(id int) error
}
