package service

import (
	"ToDoProject/internal/model"
	"ToDoProject/internal/repository"
)

type TodoStepService interface {
	GetAllStepsForTodo(todoId int) ([]model.TodoStep, error)
	CreateStep(step model.TodoStep) (model.TodoStep, error)
	UpdateStep(step model.TodoStep) (model.TodoStep, error)
	DeleteStep(id int) error
}

type todoStepService struct {
	repo repository.TodoStepRepository
}

func NewTodoStepService(repo repository.TodoStepRepository) *todoStepService {
	return &todoStepService{repo: repo}
}

func (s *todoStepService) GetAllStepsForTodo(todoId int) ([]model.TodoStep, error) {
	return s.repo.GetAllStepsForTodo(todoId)
}

func (s *todoStepService) CreateStep(step model.TodoStep) (model.TodoStep, error) {
	return s.repo.CreateStep(step)
}

func (s *todoStepService) UpdateStep(step model.TodoStep) (model.TodoStep, error) {
	return s.repo.UpdateStep(step)
}

func (s *todoStepService) DeleteStep(id int) error {
	return s.repo.DeleteStep(id)
}
