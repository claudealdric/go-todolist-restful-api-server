package models

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func NewTask(id int, title string) *Task {
	return &Task{id, title}
}

type CreateTaskDTO struct {
	Title string `json:"title"`
}

type UpdateTaskDTO struct {
	Title *string `json:"title,omitempty"`
}
