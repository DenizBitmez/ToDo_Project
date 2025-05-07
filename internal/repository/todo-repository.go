package repository

import (
	"ToDoProject/internal/model"
	"errors"
	"sync"
	"time"
)

type ToDoRepository interface {
	GetAll() ([]model.TodoList, error)
	GetAllByUsername(username string) ([]model.TodoList, error)
	GetById(id int) (*model.TodoList, error)
	Create(list model.TodoList) model.TodoList
	Update(list model.TodoList) error
	SoftDelete(id int) error
	Restore(id int) error
}

type InMemoryToDoRepository struct {
	todos  map[int]model.TodoList
	mu     sync.RWMutex
	nextID int
}

func NewInMemoryToDoRepository() *InMemoryToDoRepository {
	return &InMemoryToDoRepository{
		todos:  make(map[int]model.TodoList),
		nextID: 1,
	}
}

func (r *InMemoryToDoRepository) GetAll() ([]model.TodoList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.TodoList
	for _, todo := range r.todos {
		if !todo.IsDeleted {
			result = append(result, todo)
		}
	}
	return result, nil
}

func (r *InMemoryToDoRepository) GetAllByUsername(username string) ([]model.TodoList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.TodoList
	for _, todo := range r.todos {
		if todo.Username == username && !todo.IsDeleted {
			result = append(result, todo)
		}
	}
	return result, nil
}

func (r *InMemoryToDoRepository) GetById(id int) (*model.TodoList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists || todo.IsDeleted {
		return nil, errors.New("todo bulunamadı")
	}
	return &todo, nil
}

func (r *InMemoryToDoRepository) Create(todo model.TodoList) model.TodoList {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo.ID = r.nextID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.IsDeleted = false
	r.todos[r.nextID] = todo
	r.nextID++
	return todo
}

func (r *InMemoryToDoRepository) Update(todo model.TodoList) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingTodo, exists := r.todos[todo.ID]
	if !exists || existingTodo.IsDeleted {
		return errors.New("todo bulunamadı")
	}

	todo.UpdatedAt = time.Now()
	todo.CreatedAt = existingTodo.CreatedAt
	todo.IsDeleted = existingTodo.IsDeleted
	r.todos[todo.ID] = todo
	return nil
}

func (r *InMemoryToDoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists || todo.IsDeleted {
		return errors.New("todo bulunamadı")
	}

	now := time.Now()
	todo.DeletedAt = &now
	todo.IsDeleted = true
	todo.UpdatedAt = now
	r.todos[id] = todo
	return nil
}

func (r *InMemoryToDoRepository) Restore(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists {
		return errors.New("todo bulunamadı")
	}

	if !todo.IsDeleted {
		return errors.New("todo zaten aktif")
	}

	todo.IsDeleted = false
	todo.DeletedAt = nil
	todo.UpdatedAt = time.Now()
	r.todos[id] = todo
	return nil
}

func (r *InMemoryToDoRepository) SoftDelete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists || todo.IsDeleted {
		return errors.New("todo bulunamadı")
	}

	now := time.Now()
	todo.DeletedAt = &now
	todo.IsDeleted = true
	todo.UpdatedAt = now
	r.todos[id] = todo
	return nil
}
