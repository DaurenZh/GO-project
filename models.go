package main

import (
	"errors"
	"strings"
)

// TaskType identifies the kind of entry the application processes.
type TaskType string

const (
	TypeHabit TaskType = "habit"
	TypeNote  TaskType = "note"
	TypeTodo  TaskType = "todo"
)

// TaskState represents the processing state of a task.
type TaskState string

const (
	StatePending    TaskState = "pending"
	StateProcessing TaskState = "processing"
	StateCompleted  TaskState = "completed"
)

// Habit stores information about a habit.
type Habit struct {
	Name string
	Note
	CompletedDays int
	TotalDays     int
}

// Note stores a title and body.
type Note struct {
	Title string
	Body  string
}

// Priority represents the priority of a to-do item.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Todo stores information about a to-do item.
type Todo struct {
	Note
	Priority         Priority
	EstimatedMinutes int
}

// Task represents a unit of work processed.
type Task struct {
	ID      int
	Type    TaskType
	Payload any
	State   TaskState
	Result  any
}

// Validate checks whether a habit contains valid data.
func (h Habit) Validate() error {
	if strings.TrimSpace(h.Name) == "" {
		return errors.New("habit name cannot be empty")
	}

	if h.CompletedDays < 0 {
		return errors.New("completed days cannot be negative")
	}

	if h.TotalDays <= 0 {
		return errors.New("total days must be greater than zero")
	}

	if h.CompletedDays > h.TotalDays {
		return errors.New("completed days cannot exceed total days")
	}

	return nil
}

// Validate checks whether a note contains valid data.
func (n Note) Validate() error {
	if strings.TrimSpace(n.Title) == "" {
		return errors.New("note title cannot be empty")
	}

	if strings.TrimSpace(n.Body) == "" {
		return errors.New("note body cannot be empty")
	}

	return nil
}

// Validate checks whether a to-do item contains valid data.
func (t Todo) Validate() error {
	if strings.TrimSpace(t.Note.Title) == "" {
		return errors.New("to-do title cannot be empty")
	}

	switch t.Priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		// Valid priority.
	default:
		return errors.New("priority must be low, medium, or high")
	}

	if t.EstimatedMinutes <= 0 {
		return errors.New("estimated minutes must be greater than zero")
	}

	return nil
}
