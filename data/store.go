package data

import "github.com/claudealdric/go-todolist-restful-api-server/models"

type Store interface {
	CreateTask(task models.Task) (models.Task, error)
	DeleteTaskById(id int) error
	GetTasks() ([]models.Task, error)
}
