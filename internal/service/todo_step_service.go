package service

import (
	"ToDoProject/internal/model"
	"ToDoProject/internal/repository"
)

type TodoStepService interface {
	GetAllSteps(username string) ([]model.TodoStep, error)
	CreateStep(step model.TodoStep) (model.TodoStep, error)
	UpdateStep(step model.TodoStep) (model.TodoStep, error)
	DeleteStep(id int) error
	GetAllStepsForTodo() ([]model.TodoStep, error)
}

type todoStepService struct {
	repo repository.TodoStepRepository
}

func NewTodoStepService(repo repository.TodoStepRepository) *todoStepService {
	return &todoStepService{repo: repo}
}

func (s *todoStepService) GetAllSteps(username string) ([]model.TodoStep, error) {
	return s.repo.GetAllStepsByUsername(username)
}

func (s *todoStepService) GetAllStepsForTodo() ([]model.TodoStep, error) {
	return s.repo.GetAllSteps()
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
