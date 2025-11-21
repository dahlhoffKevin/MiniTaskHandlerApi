package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-task-api/handlers"
	"go-task-api/types"
)

func TestTaskHandler_GetAllTasks(t *testing.T) {
	store := &TaskMockStore{
		Tasks: []types.Task{
			{ID: 1, Title: "test1", Done: false},
			{ID: 2, Title: "test2", Done: false},
		},
	}

	handler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.HandleTasks)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got []types.Task
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(got))
	}

	if got[0].Title != "test1" || got[1].Title != "test2" {
		t.Fatalf("unexpected tasks: %+v", got)
	}
}

func TestTaskHandler_GetOneTask(t *testing.T) {
	store := &TaskMockStore{
		Tasks: []types.Task{
			{ID: 1, Title: "test1", Done: false},
		},
	}

	handler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/{id}", handler.HandleTasks)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got types.Task
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if got.ID != 1 || got.Title != "test1" {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	store := &TaskMockStore{}
	handler := &handlers.TaskHandler{Store: store}

	body := map[string]string{
		"title": "test1",
	}
	jsonBody, _ := json.Marshal(body)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.HandleTasks)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(jsonBody))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var got types.Task
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %+v", err)
	}

	if got.Title != "test1" {
		t.Fatalf("unexpected task in response: %+v", got)
	}

	if len(store.Created) != 1 {
		t.Fatalf("expected 1 created task in store, got %d", len(store.Created))
	}

	if store.Created[0].Title != "test1" {
		t.Fatalf("store did not record created task correctly: %+v", store.Created[0])
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	store := &TaskMockStore{
		Tasks: []types.Task{
			{ID: 1, Title: "test1", Done: false},
		},
	}

	handler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/{id}", handler.HandleTasks)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if len(store.DeletedIDs) != 1 || store.DeletedIDs[0] != 1 {
		t.Fatalf("expected DeletedIDs tp contain [1], got %v", store.DeletedIDs)
	}
}

func TestTaskHandler_MarkTaskAsDone(t *testing.T) {
	store := &TaskMockStore{
		Tasks: []types.Task{
			{ID: 1, Title: "test1", Done: false},
		},
	}

	handler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/{id}/done", handler.HandleTasks)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/1/done", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if store.Tasks[0].Done != true {
		t.Fatalf("expected task [1] to be done, task status %v", store.Tasks[0].Done)
	}
}

func TestTaskHandler_RenameTask(t *testing.T) {
	store := &TaskMockStore{
		Tasks: []types.Task{
			{ID: 1, Title: "task1", Done: false},
		},
	}

	handler := &handlers.TaskHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/{id}/rename", handler.HandleTasks)

	body := map[string]string{
		"title": "renamed task1",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/1/rename", bytes.NewReader(jsonBody))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got types.Task
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if got.Title != "renamed task1" {
		t.Fatalf("unexpected task in response: %+v", got)
	}
}
