package todo

import "github.com/tmolyakov/todo-cli/internal/task"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]task.Task, error) {
	return s.repo.List()
}

func (s *Service) Add(task task.Task) (task.Task, error) {
	return s.repo.Add(task)
}
