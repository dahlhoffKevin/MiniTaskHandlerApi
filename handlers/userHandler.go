package handlers

import (
	"encoding/json"
	"net/http"

	"go-task-api/httpError"
	"go-task-api/types"
	"go-task-api/utils"
)

type UserHandler struct {
	Store types.UserStore
}

func NewUserStore(store types.UserStore) *UserHandler {
	return &UserHandler{
		Store: store,
	}
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/users" {
			h.handleGetAllUsers(w)
			return
		}

		id := r.PathValue("id")
		if id != "" && r.URL.Path == "/users/"+id {
			h.handleGetOneUser(w, r)
			return
		}

		httpError.Write(w, httpError.New(http.StatusBadRequest, "endpoint not found"))
	case http.MethodPost:
		h.handleCreateUser(w, r)
	case http.MethodDelete:
		h.handleDeleteUser(w, r)
	default:
		httpError.Write(w, httpError.New(http.StatusMethodNotAllowed, "method not allowed"))
	}
}

func (h *UserHandler) getUserFromRequest(r *http.Request) (*types.User, *httpError.HTTPError) {
	id, err := utils.ParseIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	user, errGetById := h.Store.GetByID(id)
	if errGetById != nil {
		return nil, errGetById
	}

	return user, nil
}

func (h *UserHandler) handleGetAllUsers(w http.ResponseWriter) {
	users, err := h.Store.GetAll()
	if err != nil {
		httpError.Write(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) handleGetOneUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromRequest(r)
	if err != nil {
		httpError.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpError.Write(w, httpError.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	user, err := h.Store.Create(input.Name, input.Email)
	if err != nil {
		httpError.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, errParseID := utils.ParseIDFromRequest(r)
	if errParseID != nil {
		httpError.Write(w, errParseID)
		return
	}

	err := h.Store.Delete(id)
	if err != nil {
		httpError.Write(w, err)
		return
	}
}
