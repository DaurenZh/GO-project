package main

import (
	"errors"
	"fmt"
)

// Store manages tasks in memory.
type Store struct {
	tasks  map[int]*Task
	nextID int
}

// NewStore creates and initializes an empty task store.
func NewStore() *Store {
	return &Store{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

// Create assigns an ID to a task and stores it.
func (s *Store) Create(task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	if task.ID != 0 {
		return errors.New("task already has an ID")
	}

	if task.State != StatePending {
		return errors.New("new task must be pending")
	}

	if task.Result != nil {
		return errors.New("new task cannot have a result")
	}

	task.ID = s.nextID
	s.nextID++

	s.tasks[task.ID] = task

	return nil
}

// GetByID retrieves a task by its ID.
func (s *Store) GetByID(id int) (*Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	return task, nil
}

// List returns all stored tasks.
func (s *Store) List() []*Task {
	tasks := make([]*Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// Delete removes a task from the in-memory store.
func (s *Store) Delete(id int) error {
	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}

	delete(s.tasks, id)
	return nil
}
