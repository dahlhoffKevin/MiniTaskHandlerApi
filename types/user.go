package types

import (
	"net/http"

	"go-task-api/httpError"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"` //25414004-d988-458e-aaf9-463ec3218d8b (8-4-4-4-12)
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func GetUserFromUserID(userID uuid.UUID, users []User) (*User, *httpError.HTTPError) {
	idx, err := GetUserIndexFromUserID(userID, users)
	if err != nil {
		return nil, err
	}

	return &users[idx], nil
}

func GetUserIndexFromUserID(userID uuid.UUID, tasks []User) (int, *httpError.HTTPError) {
	// passenden index finden
	idx := -1
	for i, t := range tasks {
		if t.ID == userID {
			idx = i
			break
		}
	}

	if idx == -1 {
		return 0, httpError.New(http.StatusNotFound, "user not found")
	}

	return idx, nil
}
