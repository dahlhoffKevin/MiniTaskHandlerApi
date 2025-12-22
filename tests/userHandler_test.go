package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-task-api/handlers"
	"go-task-api/types"

	"github.com/google/uuid"
)

func TestUserHandler_GetAllUsers(t *testing.T) {
	uuidForUserOne := uuid.New()
	uuidForUserTwo := uuid.New()

	store := &UserMockStore{
		Users: []types.User{
			{ID: uuidForUserOne, Name: "Alice", Email: "alice@example.com"},
			{ID: uuidForUserTwo, Name: "Bob", Email: "bob@example.com"},
		},
	}

	handler := &handlers.UserHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.HandleUsers)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got []types.User
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}

	if got[0].Name != "Alice" || got[1].Name != "Bob" {
		t.Fatalf("unexpected users: %+v", got)
	}
}

func TestUserHandler_GetOneUser(t *testing.T) {
	uuidForUserOne := uuid.New()

	store := &UserMockStore{
		Users: []types.User{
			{ID: uuidForUserOne, Name: "Alice", Email: "alice@example.com"},
		},
	}

	handler := &handlers.UserHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/{id}", handler.HandleUsers)

	req := httptest.NewRequest(http.MethodGet, "/users/"+uuidForUserOne.String(), nil)
	rr := httptest.NewRecorder()

	// wichtig: über mux routen, NICHT direkt handler.HandleUsers
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got types.User
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if got.ID != uuidForUserOne || got.Name != "Alice" {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	store := &UserMockStore{}
	handler := &handlers.UserHandler{Store: store}

	body := map[string]string{
		"name":  "Charlie",
		"email": "charlie@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.HandleUsers)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	var got types.User
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if got.Name != "Charlie" || got.Email != "charlie@example.com" {
		t.Fatalf("unexpected user in response: %+v", got)
	}

	if len(store.Created) != 1 {
		t.Fatalf("expected 1 created user in store, got %d", len(store.Created))
	}

	if store.Created[0].Name != "Charlie" {
		t.Fatalf("store did not record created user correctly: %+v", store.Created[0])
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	uuidForUserOne := uuid.New()

	store := &UserMockStore{
		Users: []types.User{
			{ID: uuidForUserOne, Name: "Alice", Email: "alice@example.com"},
		},
	}

	handler := &handlers.UserHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/{id}", handler.HandleUsers)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+uuidForUserOne.String(), nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if len(store.DeletedIDs) != 1 || store.DeletedIDs[0] != uuidForUserOne {
		t.Fatalf("expected DeletedIDs to contain [1], got %v", store.DeletedIDs)
	}
}
