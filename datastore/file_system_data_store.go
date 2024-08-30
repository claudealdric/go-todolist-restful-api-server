package datastore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/claudealdric/go-todolist-restful-api-server/models"
)

type FileSystemDataStore struct {
	database *json.Encoder
	tasks    []models.Task
}

func NewFileSystemDataStore(file *os.File) (*FileSystemDataStore, error) {
	err := initializeDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initializing player db file, %v", err)
	}

	tasks, err := newTasks(file)

	if err != nil {
		return nil, fmt.Errorf(
			"problem loading player store from file %s, %v",
			file.Name(),
			err,
		)
	}

	return &FileSystemDataStore{
		json.NewEncoder(&tape{file}),
		tasks,
	}, nil
}

func (f *FileSystemDataStore) GetTasks() []models.Task {
	return f.tasks
}

func (f *FileSystemDataStore) CreateTask(task models.Task) models.Task {
	f.tasks = append(f.tasks, task)
	f.database.Encode(f.tasks)
	return task
}

func (f *FileSystemDataStore) DeleteTaskById(id int) {
	f.tasks = slices.DeleteFunc(f.tasks, func(task models.Task) bool {
		return task.Id == id
	})
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

func newTasks(r io.Reader) ([]models.Task, error) {
	var tasks []models.Task
	err := json.NewDecoder(r).Decode(&tasks)
	if err != nil {
		err = fmt.Errorf("problem parsing tasks, %v", err)
	}
	return tasks, err
}
