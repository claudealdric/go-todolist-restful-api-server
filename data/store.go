package data

import (
	"errors"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

var ErrResourceNotFound = errors.New("resource not found")

type Store interface {
	CreateTask(dto *models.CreateTaskDTO) (*models.Task, error)
	DeleteTaskById(id int) error
	GetTaskById(id int) (*models.Task, error)
	GetTasks() ([]models.Task, error)
	UpdateTask(task models.Task) (models.Task, error)

	CreateUser(dto *models.CreateUserDTO) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers() ([]models.User, error)
}
