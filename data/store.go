package data

import (
	"errors"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

var ErrResourceNotFound = errors.New("resource not found")

type Store interface {
	CreateTask(task models.Task) (models.Task, error)
	DeleteTaskById(id int) error
	GetTaskById(id int) (models.Task, error)
	GetTasks() ([]models.Task, error)
}
