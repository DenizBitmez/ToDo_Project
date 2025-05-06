package repository

import (
	"ToDoProject/internal/model"
	"errors"
	"sync"
	"time"
)

type ToDoRepository interface {
	GetTodosByUsername(username string) []model.TodoList
	GetById(id int) (*model.TodoList, error)
	Create(list model.TodoList) model.TodoList
	Update(list model.TodoList) error
	SoftDelete(id int) error
}

type InMemoryToDoRepository struct {
	data  []model.TodoList
	mutex sync.RWMutex
}

func NewInMemoryToDoRepository() *InMemoryToDoRepository {
	return &InMemoryToDoRepository{
		data: make([]model.TodoList, 0),
	}
}

func (repo *InMemoryToDoRepository) GetTodosByUsername(username string) []model.TodoList {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var result []model.TodoList
	for _, item := range repo.data {
		if item.DeletedAt == nil && item.Name == username {
			result = append(result, item)
		}
	}

	return result
}
func (repo *InMemoryToDoRepository) GetById(id int) (*model.TodoList, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	for _, item := range repo.data {
		if item.ID == id && item.DeletedAt == nil {
			return &item, nil
		}
	}
	return nil, errors.New("Not Found")
}

func (r *InMemoryToDoRepository) Create(list model.TodoList) model.TodoList {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	list.ID = len(r.data) + 1
	list.CreatedAt = time.Now()
	r.data = append(r.data, list)
	return list
}

func (r *InMemoryToDoRepository) Update(updated model.TodoList) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, item := range r.data {
		if item.ID == updated.ID && item.DeletedAt == nil {
			updated.UpdatedAt = time.Now()
			r.data[i] = updated
			return nil
		}
	}
	return errors.New("Not Found")
}

func (r *InMemoryToDoRepository) SoftDelete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, item := range r.data {
		if item.ID == id && item.DeletedAt == nil {
			now := time.Now()
			r.data[i].DeletedAt = &now
			return nil
		}
	}
	return errors.New("Not Found")
}
