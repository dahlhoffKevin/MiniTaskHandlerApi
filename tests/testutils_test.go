package tests

import (
	"net/http"
	"net/http/httptest"

	"go-task-api/handlers"
	"go-task-api/httpError"
	"go-task-api/types"
)

type TaskMockStore struct {
	Tasks      []types.Task
	Created    []types.Task
	DeletedIDs []int
}

func (s *TaskMockStore) GetAll() ([]types.Task, *httpError.HTTPError) {
	return s.Tasks, nil
}

func (s *TaskMockStore) GetByID(id int) (*types.Task, *httpError.HTTPError) {
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return &s.Tasks[i], nil
		}
	}
	return nil, httpError.New(http.StatusNotFound, "user not found")
}

func (s *TaskMockStore) Create(title string) types.Task {
	task := types.Task{
		ID:    len(s.Tasks) + 1,
		Title: title,
		Done:  false,
	}
	s.Tasks = append(s.Tasks, task)
	s.Created = append(s.Created, task)
	return task
}

func (s *TaskMockStore) Delete(id int) *httpError.HTTPError {
	s.DeletedIDs = append(s.DeletedIDs, id)
	return nil
}

func SetupTestHeaderForTaskTests(store *TaskMockStore) {
	handler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.HandleTasks)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
}

type UserMockStore struct {
	Users      []types.User
	Created    []types.User
	DeletedIDs []int
}

func (s *UserMockStore) GetAll() ([]types.User, *httpError.HTTPError) {
	return s.Users, nil
}

func (s *UserMockStore) GetByID(id int) (*types.User, *httpError.HTTPError) {
	for i := range s.Users {
		if s.Users[i].ID == id {
			return &s.Users[i], nil
		}
	}
	return nil, httpError.New(http.StatusNotFound, "user not found")
}

func (s *UserMockStore) Create(name string, email string) (types.User, *httpError.HTTPError) {
	user := types.User{
		ID:    len(s.Users) + 1,
		Name:  name,
		Email: email,
	}
	s.Users = append(s.Users, user)
	s.Created = append(s.Created, user)
	return user, nil
}

func (s *UserMockStore) Delete(id int) *httpError.HTTPError {
	s.DeletedIDs = append(s.DeletedIDs, id)
	return nil
}
