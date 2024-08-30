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

func (f *FileSystemDataStore) GetTasks() []models.Task {
	return f.getTasksFromFile()
}

func (f *FileSystemDataStore) CreateTask(task models.Task) models.Task {
	tasks := append(f.GetTasks(), task)
	f.overwriteFile(tasks)
	return task
}

func (f *FileSystemDataStore) DeleteTaskById(id int) {
	tasks := slices.DeleteFunc(f.GetTasks(), func(task models.Task) bool {
		return task.Id == id
	})
	f.overwriteFile(tasks)
}

func (f *FileSystemDataStore) getTasksFromFile() []models.Task {
	var tasks []models.Task
	f.decoder.Decode(&tasks)
	return tasks
}

func (f *FileSystemDataStore) overwriteFile(data any) error {
	return f.encoder.Encode(data)
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
