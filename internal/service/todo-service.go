package service

import (
	"ToDoProject/internal/model"
	"ToDoProject/internal/repository"
)

type TodoListService interface {
	GetAllByUsername(username string) ([]model.TodoList, error)
	GetAll() ([]model.TodoList, error)
	GetById(id int) (*model.TodoList, error)
	Create(todo model.TodoList) model.TodoList
	Update(todo model.TodoList) error
	Delete(todo int) error
	Restore(id int) error
}

type todoService struct {
	repo repository.ToDoRepository
}

func NewTodoService(repo repository.ToDoRepository) *todoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAll() ([]model.TodoList, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return todos, nil
}
func (s *todoService) GetAllByUsername(username string) ([]model.TodoList, error) {
	todos, err := s.repo.GetAllByUsername(username)
	if err != nil {
		return nil, err
	}
	return todos, nil
}
func (s *todoService) GetById(id int) (*model.TodoList, error) {
	return s.repo.GetById(id)
}

func (s *todoService) Create(todo model.TodoList) model.TodoList {
	return s.repo.Create(todo)
}

func (s *todoService) Update(todo model.TodoList) error {
	return s.repo.Update(todo)
}

func (s *todoService) Delete(id int) error {
	return s.repo.SoftDelete(id)
}

func (s *todoService) Restore(id int) error {
	return s.repo.Restore(id)
}
