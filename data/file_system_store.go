package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
	"github.com/claudealdric/go-todolist-restful-api-server/utils"
)

type FileSystemStore struct {
	encoder    *json.Encoder
	decoder    *json.Decoder
	lastTaskId int
	lastUserId int
}

func NewFileSystemStore(file *os.File) (*FileSystemStore, error) {
	err := initializeDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initializing player db file, %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf(
			"problem loading player store from file %s, %v",
			file.Name(),
			err,
		)
	}

	return &FileSystemStore{
		json.NewEncoder(&tape{file}),
		json.NewDecoder(&tape{file}),
		0,
		0,
	}, nil
}

func (f *FileSystemStore) GetTaskById(id int) (*models.Task, error) {
	tasks, err := f.getTasksFromFile()
	if err != nil {
		return nil, err
	}
	task, ok := utils.SliceFind(tasks, func(t models.Task) bool {
		return t.Id == id
	})
	if !ok {
		return nil, fmt.Errorf(
			"task with ID %d: %w",
			id,
			ErrResourceNotFound,
		)
	}
	return &task, nil
}

func (f *FileSystemStore) GetTasks() ([]models.Task, error) {
	tasks, err := f.getTasksFromFile()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (f *FileSystemStore) CreateTask(dto *models.CreateTaskDTO) (*models.Task, error) {
	tasks, err := f.GetTasks()
	if err != nil {
		return nil, err
	}
	newId := f.getNewTaskId()
	task := models.Task{
		Id:    newId,
		Title: dto.Title,
	}
	tasks = append(tasks, task)
	err = f.overwriteFile(tasks)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (f *FileSystemStore) DeleteTaskById(id int) error {
	tasks, err := f.GetTasks()
	if err != nil {
		return err
	}
	i := slices.IndexFunc(tasks, func(task models.Task) bool {
		return task.Id == id
	})
	if i == -1 {
		return fmt.Errorf("error with task ID %d: %w", id, ErrResourceNotFound)
	}
	tasks = slices.DeleteFunc(tasks, func(task models.Task) bool {
		return task.Id == id
	})
	return f.overwriteFile(tasks)
}

func (f *FileSystemStore) UpdateTask(task models.Task) (*models.Task, error) {
	tasks, err := f.GetTasks()
	if err != nil {
		return nil, err
	}

	taskToUpdate, err := f.GetTaskById(task.Id)
	if err != nil {
		return nil, err
	}

	taskToUpdate.Title = task.Title
	for i, t := range tasks {
		if t.Id == task.Id {
			tasks[i] = *taskToUpdate
			break
		}
	}
	err = f.overwriteFile(tasks)
	if err != nil {
		return nil, err
	}

	return taskToUpdate, nil
}

func (f *FileSystemStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	users, err := f.getUsersFromFile()
	if err != nil {
		return nil, err
	}
	user, ok := utils.SliceFind(users, func(u models.User) bool {
		return u.Email == email
	})
	if !ok {
		return nil, fmt.Errorf(
			"user with email %s: %w",
			email,
			ErrResourceNotFound,
		)
	}
	return &user, nil
}

func (f *FileSystemStore) CreateUser(dto *models.CreateUserDTO) (*models.User, error) {
	users, err := f.GetUsers()
	if err != nil {
		return nil, err
	}
	newId := f.getNewUserId()
	user := models.User{
		Id:       newId,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
	users = append(users, user)
	err = f.overwriteFile(users)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (f *FileSystemStore) GetUsers() ([]models.User, error) {
	users, err := f.getUsersFromFile()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (f *FileSystemStore) getTasksFromFile() ([]models.Task, error) {
	var tasks []models.Task
	err := f.decoder.Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %w", err)
	}
	return tasks, nil
}

func (f *FileSystemStore) getUsersFromFile() ([]models.User, error) {
	var users []models.User
	err := f.decoder.Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %w", err)
	}
	return users, nil
}

func (f *FileSystemStore) overwriteFile(data any) error {
	err := f.encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (f *FileSystemStore) getNewTaskId() int {
	newTaskId := f.lastTaskId + 1
	f.lastTaskId = newTaskId
	return newTaskId
}

func (f *FileSystemStore) getNewUserId() int {
	newUserId := f.lastUserId + 1
	f.lastUserId = newUserId
	return newUserId
}

func initializeDBFile(file *os.File) error {
	_, err := file.Seek(0, io.SeekStart)

	if err != nil {
		return fmt.Errorf(
			"problem seeking from file %s, %v",
			file.Name(),
			err,
		)
	}

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf(
			"problem getting info from file %s, %v",
			file.Name(),
			err,
		)
	}

	if info.Size() == 0 {
		_, err := file.Write([]byte("[]"))

		if err != nil {
			return fmt.Errorf(
				"problem writing to file %s, %v",
				file.Name(),
				err,
			)
		}

		_, err = file.Seek(0, io.SeekStart)

		if err != nil {
			return fmt.Errorf(
				"problem seeking from file %s, %v",
				file.Name(),
				err,
			)
		}
	}

	return nil
}
