package jsonfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tmolyakov/todo-cli/internal/task"
)

// Repository stores tasks in a JSON file.
type Repository struct {
	path string
}

func New(path string) *Repository {
	return &Repository{path: path}
}

func (r *Repository) List() ([]task.Task, error) {
	data, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return []task.Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read tasks from %q: %w", r.path, err)
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("decode tasks from %q: %w", r.path, err)
	}
	return tasks, nil
}

func (r *Repository) Add(t task.Task) (task.Task, error) {
	tasks, err := r.List()
	if err != nil {
		return task.Task{}, fmt.Errorf("get list %w", err)
	}

	var maxID uint64
	for _, existing := range tasks {
		if existing.ID > maxID {
			maxID = existing.ID
		}
	}
	t.ID = maxID + 1
	err = t.Validate()
	if err != nil {
		return task.Task{}, fmt.Errorf("task validation err %w", err)
	}
	tasks = append(tasks, t)

	data, err := json.Marshal(tasks)
	if err != nil {
		return task.Task{}, fmt.Errorf("encode task err %w", err)
	}

	err = os.WriteFile(r.path, data, 0644)
	if err != nil {
		return task.Task{}, fmt.Errorf("write task to file err %w", err)
	}
	return t, nil
}
