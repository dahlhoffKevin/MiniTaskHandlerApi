package types

import (
	"go-task-api/httpError"

	"github.com/google/uuid"
)

type UserStore interface {
	GetAll() ([]User, *httpError.HTTPError)
	GetByID(id uuid.UUID) (*User, *httpError.HTTPError)
	Create(name string, email string) (User, *httpError.HTTPError)
	Delete(id uuid.UUID) *httpError.HTTPError
}
