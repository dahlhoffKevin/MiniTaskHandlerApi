package storage

import (
	"go-task-api/httpError"
	"go-task-api/types"
)

type InMemoryUserStore struct {
	Users  []types.User
	NextID int
}

func (memoryUserStore *InMemoryUserStore) GetAll() ([]types.User, *httpError.HTTPError) {
	if len(memoryUserStore.Users) == 0 {
		return nil, httpError.New(404, "no users found")
	}

	return memoryUserStore.Users, nil
}

func (memoryUserStore *InMemoryUserStore) GetByID(id int) (*types.User, *httpError.HTTPError) {
	if len(memoryUserStore.Users) == 0 {
		return nil, httpError.New(404, "no users found")
	}

	user, err := types.GetUserFromUserID(id, memoryUserStore.Users)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (memoryUserStore *InMemoryUserStore) Create(name string, email string) (types.User, *httpError.HTTPError) {
	if name == "" || email == "" {
		return types.User{}, httpError.New(400, "name or email missing")
	}

	user := types.User{
		ID:    memoryUserStore.NextID,
		Name:  name,
		Email: email,
	}
	memoryUserStore.NextID++
	memoryUserStore.Users = append(memoryUserStore.Users, user)

	return user, nil
}

func (memoryUserStore *InMemoryUserStore) Delete(id int) *httpError.HTTPError {
	if len(memoryUserStore.Users) == 0 {
		return httpError.New(404, "no users found")
	}

	idx, err := types.GetUserIndexFromUserID(id, memoryUserStore.Users)
	if err != nil {
		return err
	}

	memoryUserStore.Users = append(memoryUserStore.Users[:idx], memoryUserStore.Users[idx+1:]...)
	return nil
}
