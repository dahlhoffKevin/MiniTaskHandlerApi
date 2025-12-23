package types

import (
	"net/http"

	"go-task-api/httpError"

	"github.com/google/uuid"
)

type Task struct {
	ID     uuid.UUID
	Title  string
	Done   bool
	UserID uuid.UUID
}

func (t *Task) MarkDone() {
	t.Done = true
}

func (t *Task) Rename(newTitle string) {
	t.Title = newTitle
}

func GetTaskFromTaskID(taskID uuid.UUID, tasks []Task) (*Task, *httpError.HTTPError) {
	idx, err := GetTaskIndexFromTaskID(taskID, tasks)
	if err != nil {
		return nil, err
	}

	return &tasks[idx], nil
}

func GetTaskIndexFromTaskID(taskID uuid.UUID, tasks []Task) (int, *httpError.HTTPError) {
	// passenden Index finden
	idx := -1
	for i, t := range tasks {
		if t.ID == taskID {
			idx = i
			break
		}
	}

	if idx == -1 {
		return 0, httpError.New(http.StatusNotFound, "task not found")
	}

	return idx, nil
}
