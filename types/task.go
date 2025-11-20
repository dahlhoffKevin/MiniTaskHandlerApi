package types

import (
	"go-task-api/httpError"
	"net/http"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

func (t *Task) MarkDone() {
	t.Done = true
}

func (t *Task) Rename(newTitle string) {
	t.Title = newTitle
}

func GetTaskFromTaskID(taskID int, tasks []Task) (*Task, *httpError.HTTPError) {
	idx, err := GetTaskIndexFromTaskID(taskID, tasks)
	if err != nil {
		return nil, err
	}

	return &tasks[idx], nil
}

func GetTaskIndexFromTaskID(taskID int, tasks []Task) (int, *httpError.HTTPError) {
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
