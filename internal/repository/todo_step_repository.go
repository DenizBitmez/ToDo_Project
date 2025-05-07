package repository

import (
	"ToDoProject/internal/model"
	"errors"
	"sync"
	"time"
)

type TodoStepRepository interface {
	GetAllSteps() ([]model.TodoStep, error)
	CreateStep(step model.TodoStep) (model.TodoStep, error)
	UpdateStep(step model.TodoStep) (model.TodoStep, error)
	DeleteStep(id int) error
	GetAllStepsByUsername(username string) ([]model.TodoStep, error)
}

type InMemoryTodoStepRepository struct {
	data  []model.TodoStep
	mutex sync.RWMutex
}

func NewInMemoryTodoStepRepository() *InMemoryTodoStepRepository {
	return &InMemoryTodoStepRepository{
		data: make([]model.TodoStep, 0),
	}
}
func (repo *InMemoryTodoStepRepository) GetAllSteps() ([]model.TodoStep, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var steps []model.TodoStep
	for _, step := range repo.data {
		if step.DeletedAt == nil {
			steps = append(steps, step)
		}
	}
	return steps, nil
}

func (repo *InMemoryTodoStepRepository) GetAllStepsByUsername(username string) ([]model.TodoStep, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var steps []model.TodoStep
	for _, step := range repo.data {
		if step.Username == username && step.DeletedAt == nil {
			steps = append(steps, step)
		}
	}
	return steps, nil
}

func (repo *InMemoryTodoStepRepository) CreateStep(step model.TodoStep) (model.TodoStep, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	step.ID = len(repo.data) + 1
	repo.data = append(repo.data, step)

	return step, nil
}

func (repo *InMemoryTodoStepRepository) UpdateStep(step model.TodoStep) (model.TodoStep, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, s := range repo.data {
		if s.ID == step.ID {
			repo.data[i] = step
			return step, nil
		}
	}

	return model.TodoStep{}, errors.New("step not found")
}

func (repo *InMemoryTodoStepRepository) DeleteStep(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, s := range repo.data {
		if s.ID == id {
			now := time.Now()
			repo.data[i].DeletedAt = &now
			return nil
		}
	}

	return errors.New("step not found")
}
