package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Status describes whether a task is still active or completed.
type Status string

const (
	StatusActive Status = "active"
	StatusDone   Status = "done"
)

func (s Status) Valid() bool {
	return s == StatusActive || s == StatusDone
}

// Priority describes how important a task is.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) Valid() bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}

// Task is a todo item. Due is optional and uses the YYYY-MM-DD format.
type Task struct {
	ID       uint64   `json:"id"`
	Text     string   `json:"text"`
	Status   Status   `json:"status"`
	Priority Priority `json:"priority"`
	Due      string   `json:"due,omitempty"`
}

// Validate checks that a task is ready to be stored.
func (t Task) Validate() error {
	if t.ID == 0 {
		return errors.New("task ID must be positive")
	}
	if strings.TrimSpace(t.Text) == "" {
		return errors.New("task text must not be empty")
	}
	if !t.Status.Valid() {
		return fmt.Errorf("invalid task status %q", t.Status)
	}
	if !t.Priority.Valid() {
		return fmt.Errorf("invalid task priority %q", t.Priority)
	}
	if err := ValidateDueDate(t.Due); err != nil {
		return err
	}
	return nil
}

// ValidateDueDate accepts an empty value or a real calendar date in YYYY-MM-DD form.
func ValidateDueDate(value string) error {
	if value == "" {
		return nil
	}
	date, err := time.Parse("2006-01-02", value)
	if err != nil || date.Format("2006-01-02") != value {
		return fmt.Errorf("invalid due date %q: expected YYYY-MM-DD", value)
	}
	return nil
}
