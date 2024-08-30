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
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileSystemDataStore(file *os.File) (*FileSystemDataStore, error) {
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

	return &FileSystemDataStore{
		json.NewEncoder(&tape{file}),
		json.NewDecoder(&tape{file}),
	}, nil
}

func (f *FileSystemDataStore) GetTasks() ([]models.Task, error) {
	tasks, err := f.getTasksFromFile()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (f *FileSystemDataStore) CreateTask(task models.Task) (models.Task, error) {
	tasks, err := f.GetTasks()
	if err != nil {
		return models.Task{}, err
	}
	tasks = append(tasks, task)
	err = f.overwriteFile(tasks)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (f *FileSystemDataStore) DeleteTaskById(id int) error {
	tasks, err := f.GetTasks()
	if err != nil {
		return err
	}
	tasks = slices.DeleteFunc(tasks, func(task models.Task) bool {
		return task.Id == id
	})
	return f.overwriteFile(tasks)
}

func (f *FileSystemDataStore) getTasksFromFile() ([]models.Task, error) {
	var tasks []models.Task
	err := f.decoder.Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("error reading the file: %w", err)
	}
	return tasks, nil
}

func (f *FileSystemDataStore) overwriteFile(data any) error {
	err := f.encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
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
