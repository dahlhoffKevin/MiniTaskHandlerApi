package storage

import (
	"go-task-api/httpError"
	"go-task-api/types"

	"github.com/google/uuid"
)

type InMemoryTaskStore struct {
	Tasks  []types.Task
	NextID int
}

func (memoryTaskStore *InMemoryTaskStore) GetAll() ([]types.Task, *httpError.HTTPError) {
	if len(memoryTaskStore.Tasks) == 0 {
		return nil, httpError.New(404, "no tasks found")
	}

	return memoryTaskStore.Tasks, nil
}

func (memoryTaskStore *InMemoryTaskStore) GetByID(id uuid.UUID) (*types.Task, *httpError.HTTPError) {
	if len(memoryTaskStore.Tasks) == 0 {
		return nil, httpError.New(404, "no tasks found")
	}

	task, err := types.GetTaskFromTaskID(id, memoryTaskStore.Tasks)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (memoryTaskStore *InMemoryTaskStore) Create(title string) types.Task {
	task := types.Task{
		ID:    uuid.New(),
		Title: title,
		Done:  false,
	}
	memoryTaskStore.NextID++
	memoryTaskStore.Tasks = append(memoryTaskStore.Tasks, task)

	return task
}

func (memoryTaskStore *InMemoryTaskStore) Delete(id uuid.UUID) *httpError.HTTPError {
	if len(memoryTaskStore.Tasks) == 0 {
		return httpError.New(404, "no tasks found")
	}

	idx, err := types.GetTaskIndexFromTaskID(id, memoryTaskStore.Tasks)
	if err != nil {
		return err
	}

	memoryTaskStore.Tasks = append(memoryTaskStore.Tasks[:idx], memoryTaskStore.Tasks[idx+1:]...)
	return nil
}
