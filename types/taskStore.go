package types

import (
	"go-task-api/httpError"
)

type TaskStore interface {
	GetAll() ([]Task, *httpError.HTTPError)
	GetByID(id int) (*Task, *httpError.HTTPError)
	Create(title string) Task
	Delete(id int) *httpError.HTTPError
}
