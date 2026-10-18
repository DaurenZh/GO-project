package main

import (
	"strings"
	"testing"
)

// newTestService creates an empty service for each test.
func newTestService() *Service {
	return NewService(NewStore())
}

// TestTaskCreation checks that valid tasks are created correctly.
func TestTaskCreation(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeHabit, Habit{
		Name:          "Exercise",
		Note:          Note{Title: "Exercise", Body: "Exercise regularly"},
		CompletedDays: 5,
		TotalDays:     7,
	})

	if err != nil {
		t.Fatalf("expected task creation to succeed, got error: %v", err)
	}

	if task.ID == 0 {
		t.Fatal("expected task to receive a non-zero ID")
	}

	if task.Type != TypeHabit {
		t.Errorf("expected type %q, got %q", TypeHabit, task.Type)
	}

	if task.State != StatePending {
		t.Errorf("expected state %q, got %q", StatePending, task.State)
	}

	if task.Result != nil {
		t.Error("expected a newly created task to have no result")
	}
}

// TestInvalidInput checks that invalid payloads are rejected.
func TestInvalidInput(t *testing.T) {
	service := newTestService()

	tests := []struct {
		name     string
		taskType TaskType
		payload  any
	}{
		{
			name:     "empty habit name",
			taskType: TypeHabit,
			payload: Habit{
				Name:          " ",
				CompletedDays: 1,
				TotalDays:     7,
			},
		},
		{
			name:     "empty note title",
			taskType: TypeNote,
			payload: Note{
				Title: " ",
				Body:  "Some content",
			},
		},
		{
			name:     "invalid to-do priority",
			taskType: TypeTodo,
			payload: Todo{
				Note:             Note{Title: "Study Go", Body: "Practice"},
				Priority:         "critical",
				EstimatedMinutes: 30,
			},
		},
		{
			name:     "wrong payload type",
			taskType: TypeHabit,
			payload:  "this is not a Habit",
		},
		{
			name:     "unsupported task type",
			taskType: "unknown",
			payload:  "some payload",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task, err := service.CreateTask(test.taskType, test.payload)

			if err == nil {
				t.Fatal("expected an error for invalid input")
			}

			if task != nil {
				t.Error("expected no task to be returned for invalid input")
			}
		})
	}
}

// TestTaskRetrieval checks that a stored task can be retrieved by ID.
func TestTaskRetrieval(t *testing.T) {
	service := newTestService()

	created, err := service.CreateTask(TypeNote, Note{
		Title: "My Note",
		Body:  "Remember to study Go",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	retrieved, err := service.GetTask(created.ID)
	if err != nil {
		t.Fatalf("expected retrieval to succeed, got error: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, retrieved.ID)
	}
}

// TestMissingTask checks that retrieving a nonexistent task returns an error.
func TestMissingTask(t *testing.T) {
	service := newTestService()

	_, err := service.GetTask(999)

	if err == nil {
		t.Fatal("expected an error for a missing task")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected a not-found error, got: %v", err)
	}
}

// TestTaskProcessing checks that processing completes a pending task.
func TestTaskProcessing(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeHabit, Habit{
		Name:          "Reading",
		Note:          Note{Title: "Reading", Body: "Read every day"},
		CompletedDays: 3,
		TotalDays:     5,
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	processed, err := service.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("expected processing to succeed, got error: %v", err)
	}

	if processed.State != StateCompleted {
		t.Errorf("expected state %q, got %q", StateCompleted, processed.State)
	}

	if processed.Result == nil {
		t.Fatal("expected processing to store a result")
	}
}

// TestResultCorrectness checks that processing calculates the correct result.
func TestResultCorrectness(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeHabit, Habit{
		Name:          "Exercise",
		Note:          Note{Title: "Exercise", Body: "Exercise regularly"},
		CompletedDays: 3,
		TotalDays:     4,
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	processed, err := service.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("process task: %v", err)
	}

	result, ok := processed.Result.(HabitResult)
	if !ok {
		t.Fatalf("expected HabitResult, got %T", processed.Result)
	}

	expectedPercent := 75.0
	if result.CompletionPercent != expectedPercent {
		t.Errorf(
			"expected completion percentage %.1f, got %.1f",
			expectedPercent,
			result.CompletionPercent,
		)
	}

	if result.CompletedDays != 3 || result.TotalDays != 4 {
		t.Errorf("unexpected habit statistics: %+v", result)
	}
}

// TestTaskState checks the task lifecycle from pending to completed.
func TestTaskState(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Learning",
		Body:  "Learning Go testing",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if task.State != StatePending {
		t.Fatalf("new task should be pending, got %q", task.State)
	}

	processed, err := service.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("process task: %v", err)
	}

	if processed.State != StateCompleted {
		t.Errorf("processed task should be completed, got %q", processed.State)
	}
}

// TestInvalidOperation checks that a completed task cannot be processed again.
func TestInvalidOperation(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Testing",
		Body:  "Check task state",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if _, err := service.ProcessTask(task.ID); err != nil {
		t.Fatalf("first processing should succeed: %v", err)
	}

	_, err = service.ProcessTask(task.ID)
	if err == nil {
		t.Fatal("expected an error when processing a completed task again")
	}
}

// TestMultipleTasks checks that tasks have independent IDs and results.
func TestMultipleTasks(t *testing.T) {
	service := newTestService()

	first, err := service.CreateTask(TypeNote, Note{
		Title: "First Note",
		Body:  "First body",
	})
	if err != nil {
		t.Fatalf("create first task: %v", err)
	}

	second, err := service.CreateTask(TypeTodo, Todo{
		Note:             Note{Title: "Second Task", Body: "Second body"},
		Priority:         PriorityHigh,
		EstimatedMinutes: 45,
	})
	if err != nil {
		t.Fatalf("create second task: %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("expected unique IDs, but both tasks have ID %d", first.ID)
	}

	if _, err := service.ProcessTask(first.ID); err != nil {
		t.Fatalf("process first task: %v", err)
	}

	if _, err := service.ProcessTask(second.ID); err != nil {
		t.Fatalf("process second task: %v", err)
	}

	firstStored, err := service.GetTask(first.ID)
	if err != nil {
		t.Fatalf("retrieve first task: %v", err)
	}

	secondStored, err := service.GetTask(second.ID)
	if err != nil {
		t.Fatalf("retrieve second task: %v", err)
	}

	if firstStored.Type != TypeNote {
		t.Errorf("expected first task type %q, got %q", TypeNote, firstStored.Type)
	}

	if secondStored.Type != TypeTodo {
		t.Errorf("expected second task type %q, got %q", TypeTodo, secondStored.Type)
	}

	if firstStored.Result == nil || secondStored.Result == nil {
		t.Error("expected both tasks to have independent processing results")
	}
}

func TestUpdatePendingTask(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Old title",
		Body:  "Old body",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	originalID := task.ID

	updated, err := service.UpdateTask(task.ID, Note{
		Title: "New title",
		Body:  "New body",
	})
	if err != nil {
		t.Fatalf("update task: %v", err)
	}

	if updated.ID != originalID {
		t.Errorf("expected ID %d, got %d", originalID, updated.ID)
	}

	note, ok := updated.Payload.(Note)
	if !ok {
		t.Fatalf("expected Note payload, got %T", updated.Payload)
	}

	if note.Title != "New title" {
		t.Errorf("expected updated title, got %q", note.Title)
	}

	if note.Body != "New body" {
		t.Errorf("expected updated body, got %q", note.Body)
	}

	if updated.State != StatePending {
		t.Errorf("expected pending state, got %q", updated.State)
	}
}

func TestUpdateRejectsInvalidPayload(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Original",
		Body:  "Original body",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	_, err = service.UpdateTask(task.ID, Note{
		Title: " ",
		Body:  "Updated body",
	})

	if err == nil {
		t.Fatal("expected an error for an invalid updated note")
	}

	stored, err := service.GetTask(task.ID)
	if err != nil {
		t.Fatalf("retrieve task: %v", err)
	}

	note := stored.Payload.(Note)
	if note.Title != "Original" {
		t.Errorf("invalid update changed the title to %q", note.Title)
	}
}

func TestUpdateRejectsWrongPayloadType(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Original",
		Body:  "Original body",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	_, err = service.UpdateTask(task.ID, Habit{
		Name:          "Exercise",
		CompletedDays: 1,
		TotalDays:     7,
	})

	if err == nil {
		t.Fatal("expected an error when payload type does not match task type")
	}
}

func TestUpdateRejectsCompletedTask(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Original",
		Body:  "Original body",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if _, err := service.ProcessTask(task.ID); err != nil {
		t.Fatalf("process task: %v", err)
	}

	_, err = service.UpdateTask(task.ID, Note{
		Title: "Changed",
		Body:  "Changed body",
	})

	if err == nil {
		t.Fatal("expected an error when updating a completed task")
	}
}

func TestDeletePendingTask(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Temporary note",
		Body:  "This note will be deleted",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if err := service.DeleteTask(task.ID); err != nil {
		t.Fatalf("delete task: %v", err)
	}

	_, err = service.GetTask(task.ID)
	if err == nil {
		t.Fatal("expected an error when retrieving a deleted task")
	}
}

func TestDeleteCompletedTask(t *testing.T) {
	service := newTestService()

	task, err := service.CreateTask(TypeNote, Note{
		Title: "Completed note",
		Body:  "This note has been processed",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	if _, err := service.ProcessTask(task.ID); err != nil {
		t.Fatalf("process task: %v", err)
	}

	err = service.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("delete completed task: %v", err)
	}

	if _, err := service.GetTask(task.ID); err == nil {
		t.Fatal("expected deleted completed task to be unavailable")
	}
}

func TestDeleteMissingTask(t *testing.T) {
	service := newTestService()

	err := service.DeleteTask(999)
	if err == nil {
		t.Fatal("expected an error when deleting a nonexistent task")
	}
}
