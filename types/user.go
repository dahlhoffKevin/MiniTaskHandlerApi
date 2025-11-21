package types

import (
	"go-task-api/httpError"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetUserFromUserID(userID int, users []User) (*User, *httpError.HTTPError) {
	idx, err := GetUserIndexFromUserID(userID, users)
	if err != nil {
		return nil, err
	}

	return &users[idx], nil
}

func GetUserIndexFromUserID(userID int, tasks []User) (int, *httpError.HTTPError) {
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
