package types

import (
	"go-task-api/httpError"
)

type UserStore interface {
	GetAll() ([]User, *httpError.HTTPError)
	GetByID(id int) (*User, *httpError.HTTPError)
	Create(name string, email string) (User, *httpError.HTTPError)
	Delete(id int) *httpError.HTTPError
}
