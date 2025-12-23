package types

import (
	"go-task-api/httpError"

	"github.com/google/uuid"
)

type TaskStore interface {
	GetAll() ([]Task, *httpError.HTTPError)
	GetByID(id uuid.UUID) (*Task, *httpError.HTTPError)
	Create(title string, userid uuid.UUID) (Task, *httpError.HTTPError)
	Delete(id uuid.UUID) *httpError.HTTPError
}
