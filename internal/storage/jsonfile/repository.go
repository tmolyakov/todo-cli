package jsonfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	tempFile, err := os.CreateTemp(filepath.Dir(r.path), ".tmp-*")
	if err != nil {
		return task.Task{}, fmt.Errorf("create temp file err %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return task.Task{}, fmt.Errorf("write temp file err %w", err)
	}
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return task.Task{}, fmt.Errorf("sync temp file err %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return task.Task{}, fmt.Errorf("close temp file err %w", err)
	}
	if err := os.Rename(tempFile.Name(), r.path); err != nil {
		return task.Task{}, fmt.Errorf("rename temp file err %w", err)
	}

	return t, nil
}
