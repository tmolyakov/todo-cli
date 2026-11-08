package todo

import "github.com/tmolyakov/todo-cli/internal/task"

type Repository interface {
	List() ([]task.Task, error)
	Add(task task.Task) (task.Task, error)
}
