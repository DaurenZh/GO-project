package main

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Processor defines the behavior required to process a domain object.
type Processor interface {
	Process() (any, error)
}

// HabitResult stores the result of habit analysis.
type HabitResult struct {
	Name              string
	CompletedDays     int
	TotalDays         int
	CompletionPercent float64
}

// NoteResult stores the result of text analysis.
type NoteResult struct {
	Title          string
	WordCount      int
	CharacterCount int
}

// TodoResult stores the result of to-do analysis.
type TodoResult struct {
	Title            string
	Priority         Priority
	PriorityCategory string
}

// Process calculates the completion statistics for a habit.
func (h Habit) Process() (any, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}

	completionPercent := float64(h.CompletedDays) /
		float64(h.TotalDays) * 100

	return HabitResult{
		Name:              h.Name,
		CompletedDays:     h.CompletedDays,
		TotalDays:         h.TotalDays,
		CompletionPercent: completionPercent,
	}, nil
}

// Process analyzes the text of a note.
func (n Note) Process() (any, error) {
	if err := n.Validate(); err != nil {
		return nil, err
	}

	words := strings.Fields(n.Body)
	characters := utf8.RuneCountInString(n.Body)

	return NoteResult{
		Title:          n.Title,
		WordCount:      len(words),
		CharacterCount: characters,
	}, nil
}

// Process evaluates the priority of a to-do item.
func (t Todo) Process() (any, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	var category string

	switch t.Priority {
	case PriorityHigh:
		category = "urgent"
	case PriorityMedium:
		category = "normal"
	case PriorityLow:
		category = "whenever possible"
	default:
		return nil, errors.New("unsupported priority")
	}

	return TodoResult{
		Title:            t.Note.Title,
		Priority:         t.Priority,
		PriorityCategory: category,
	}, nil
}
