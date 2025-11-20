package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-task-api/httpError"
	"go-task-api/types"
	"go-task-api/utils"
)

type TaskHandler struct {
	Store types.TaskStore
}

func NewTaskHandler(store types.TaskStore) *TaskHandler {
	return &TaskHandler{
		Store: store,
	}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/tasks" {
			h.handleGetAllTasks(w)
			return
		}

		id := r.PathValue("id")
		if id != "" && r.URL.Path == "/tasks/"+id {
			h.handleGetOneTask(w, r)
		}

	case http.MethodPost:
		h.handleCreateTask(w, r)
	case http.MethodDelete:
		h.handleDeleteTask(w, r)
	case http.MethodPatch:
		id := r.PathValue("id")
		if id == "" {
			httpError.Write(w, httpError.New(http.StatusBadRequest, "id value cannot be null"))
			return
		}

		switch {
		case strings.HasSuffix(r.URL.Path, "/done"):
			h.handleTaskMarkAsDone(w, r)
		case strings.HasSuffix(r.URL.Path, "/rename"):
			h.handleTaskRenameTitle(w, r)
		default:
			httpError.Write(w, httpError.New(http.StatusNotFound, "endpoint method not found"))
			return
		}
	default:
		httpError.Write(w, httpError.New(http.StatusMethodNotAllowed, "method not allowed"))
	}
}

func (h *TaskHandler) getTaskFromRequest(r *http.Request) (*types.Task, *httpError.HTTPError) {
	id, err := utils.ParseIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	task, errGetByID := h.Store.GetByID(id)
	if errGetByID != nil {
		return nil, errGetByID
	}

	return task, nil
}

func (h *TaskHandler) handleGetAllTasks(w http.ResponseWriter) {
	tasks, err := h.Store.GetAll()
	if err != nil {
		httpError.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) handleGetOneTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		httpError.Write(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpError.Write(w, httpError.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	task := h.Store.Create(input.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) handleTaskMarkAsDone(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		httpError.Write(w, err)
		return
	}

	task.MarkDone()
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) handleTaskRenameTitle(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpError.Write(w, httpError.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	task, err := h.getTaskFromRequest(r)
	if err != nil {
		httpError.Write(w, err)
		return
	}

	task.Rename(input.Title)
}
