package main

import (
	"errors"
	"fmt"
)

// Service coordinates task creation and processing.
type Service struct {
	store *Store
}

// NewService creates a service using the provided store.
func NewService(store *Store) *Service {
	return &Service{
		store: store,
	}
}

// CreateTask validates and stores a new task.
func (s *Service) CreateTask(taskType TaskType, payload any) (*Task, error) {
	if s.store == nil {
		return nil, errors.New("task store is not initialized")
	}

	var processor Processor

	switch taskType {
	case TypeHabit:
		value, ok := payload.(Habit)
		if !ok {
			return nil, errors.New("habit task requires a Habit payload")
		}
		processor = value

	case TypeNote:
		value, ok := payload.(Note)
		if !ok {
			return nil, errors.New("note task requires a Note payload")
		}
		processor = value

	case TypeTodo:
		value, ok := payload.(Todo)
		if !ok {
			return nil, errors.New("to-do task requires a Todo payload")
		}
		processor = value

	default:
		return nil, fmt.Errorf("unsupported task type: %q", taskType)
	}

	switch value := processor.(type) {
	case Habit:
		if err := value.Validate(); err != nil {
			return nil, err
		}
	case Note:
		if err := value.Validate(); err != nil {
			return nil, err
		}
	case Todo:
		if err := value.Validate(); err != nil {
			return nil, err
		}
	}

	task := &Task{
		Type:    taskType,
		Payload: payload,
		State:   StatePending,
	}

	if err := s.store.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetAllTasks retrieves all tasks from the store.
func (s *Service) GetAllTasks() ([]*Task, error) {
	if s.store == nil {
		return nil, errors.New("task store is not initialized")
	}

	return s.store.List(), nil
}

// GetTask retrieves a task by ID.
func (s *Service) GetTask(id int) (*Task, error) {
	if s.store == nil {
		return nil, errors.New("task store is not initialized")
	}

	return s.store.GetByID(id)
}

// ProcessTask processes a pending task and stores its result.
func (s *Service) ProcessTask(id int) (*Task, error) {
	task, err := s.GetTask(id)
	if err != nil {
		return nil, err
	}

	if task.State != StatePending {
		return nil, fmt.Errorf(
			"cannot process task %d: state is %q",
			id,
			task.State,
		)
	}

	processor, ok := task.Payload.(Processor)
	if !ok {
		return nil, errors.New("task payload does not implement Processor")
	}

	task.State = StateProcessing

	result, err := processor.Process()
	if err != nil {
		task.State = StatePending
		task.Result = nil
		return nil, fmt.Errorf("process task %d: %w", id, err)
	}

	task.Result = result
	task.State = StateCompleted

	return task, nil

}

// UpdateTask replaces a pending task's payload after validating it.
func (s *Service) UpdateTask(id int, payload any) (*Task, error) {
	task, err := s.GetTask(id)
	if err != nil {
		return nil, err
	}

	if task.State != StatePending {
		return nil, fmt.Errorf(
			"cannot update task %d: state is %q",
			id,
			task.State,
		)
	}

	var processor Processor

	switch task.Type {
	case TypeHabit:
		value, ok := payload.(Habit)
		if !ok {
			return nil, errors.New("habit task requires a Habit payload")
		}
		processor = value

	case TypeNote:
		value, ok := payload.(Note)
		if !ok {
			return nil, errors.New("note task requires a Note payload")
		}
		processor = value

	case TypeTodo:
		value, ok := payload.(Todo)
		if !ok {
			return nil, errors.New("to-do task requires a Todo payload")
		}
		processor = value

	default:
		return nil, fmt.Errorf("unsupported task type: %q", task.Type)
	}

	switch value := processor.(type) {
	case Habit:
		if err := value.Validate(); err != nil {
			return nil, err
		}
	case Note:
		if err := value.Validate(); err != nil {
			return nil, err
		}
	case Todo:
		if err := value.Validate(); err != nil {
			return nil, err
		}
	}

	task.Payload = payload

	return task, nil
}

// DeleteTask removes a task from the store.
func (s *Service) DeleteTask(id int) error {
	if s.store == nil {
		return errors.New("task store is not initialized")
	}

	_, err := s.store.GetByID(id)
	if err != nil {
		return err
	}

	return s.store.Delete(id)
}
