# NotesToDo

## Overview

LifeTrack is a command-line application written in Go for managing habits, notes, and to-do items.

The application demonstrates Go structures, methods, interfaces, pointers, maps, slices, error handling, in-memory storage, and unit testing.

## Features

- Create habit, note, and to-do tasks.
- Validate input before storing tasks.
- Assign a unique ID to every task.
- Store tasks in an in-memory map.
- Retrieve tasks by ID and report errors for missing tasks.
- Process different task types using a shared interface.
- Store processing results in their corresponding tasks.
- Reject unsupported task types, invalid payloads, and invalid state transitions.
- Test application behavior using Go's built-in testing package.
- Update pending tasks after validating the new payload.
- Delete tasks from in-memory storage.
- Reject updates of completed tasks.

## Task Types and Processing

### Habit

A habit contains a name, a note, the number of completed days, and the total number of tracked days.

Processing calculates the completion percentage:

`completion percentage = completed days / total days * 100`

### Note

A note contains a title and body.

Processing calculates the word count, character count, and estimated reading time. The application assumes an average reading speed of 200 words per minute, with a minimum estimate of one minute for a non-empty note.

### To-do

A to-do item contains a title, body, priority, and estimated completion time.

Processing categorizes its priority:
- High: urgent
- Medium: normal
- Low: whenever possible

## Architecture

The application separates its responsibilities across several files:

- `main.go` — demonstrates application usage.
- `models.go` — defines task types, states, domain structures, and validation.
- `processor.go` — defines the `Processor` interface and processing methods.
- `store.go` — stores and retrieves tasks in memory.
- `service.go` — validates, creates, retrieves, and processes tasks.
- `service_test.go` — tests the application's behavior.

The `Processor` interface allows habits, notes, and to-do items to provide their own processing logic through a common method.

The `Store` uses a map keyed by task ID. A `Service` coordinates validation, storage, and processing.

## Task Lifecycle

Every new task starts in the `pending` state.

The normal lifecycle is:

`Pending → Processing → Completed`

1. **Pending:** The task has been created, validated, and stored, but has not been processed.
2. **Processing:** The service begins executing the task's processing method.
3. **Completed:** Processing succeeds, and the result is stored in the task.

If processing fails, the service clears the result and returns the task to `pending`, allowing a later retry.

Processing a task that is not `pending` is rejected. Therefore, a completed task cannot be processed again.

The application processes tasks synchronously. The `processing` state is an intermediate state during a processing call rather than a state exposed between separate commands.

## Storage

Tasks are stored in memory using a Go map.

Each task receives an incrementing integer ID when created. Tasks remain available while the application is running, but they are not persisted after the process exits.

## Validation and Error Handling

The application rejects:
- Empty habit names.
- Invalid habit day counts.
- Empty note titles or bodies.
- Unsupported to-do priorities.
- Invalid estimated completion times.
- Payloads that do not match the specified task type.
- Unsupported task types.
- Attempts to retrieve nonexistent tasks.
- Attempts to process tasks in an invalid state.

Errors are returned to the caller rather than silently ignored.

## Requirements

- Go 1.27 or later.

## Running the Application

Run the application from the project directory:

```bash
go run .
```

## Running Tests

Run all tests:

```bash
go test ./...
```

Display each test:

```bash
go test -v ./...
```

Run static analysis:

```bash
go vet ./...
```

Check formatting:

```bash
gofmt -l .
```

The test suite covers task creation, invalid input, retrieval, missing-task errors, processing, result correctness, lifecycle states, invalid operations, and multiple independent tasks.