package main

import (
	"fmt"
	"log"
)

func main() {
	store := NewStore()
	service := NewService(store)

	fmt.Println("=== LifeTrack ===")

	habit, err := service.CreateTask(TypeHabit, Habit{
		Name: "Exercise",
		Note: Note{
			Title: "Exercise",
			Body:  "Exercise regularly",
		},
		CompletedDays: 5,
		TotalDays:     7,
	})
	if err != nil {
		log.Fatal(err)
	}

	note, err := service.CreateTask(TypeNote, Note{
		Title: "Go Learning",
		Body:  "Today I learned about structs, methods, and interfaces.",
	})
	if err != nil {
		log.Fatal(err)
	}

	todo, err := service.CreateTask(TypeTodo, Todo{
		Note: Note{
			Title: "Complete Go assignment",
			Body:  "Implement processing and tests",
		},
		Priority:         PriorityHigh,
		EstimatedMinutes: 60,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created habit: ID=%d, state=%s\n", habit.ID, habit.State)
	fmt.Printf("Created note: ID=%d, state=%s\n", note.ID, note.State)
	fmt.Printf("Created to-do: ID=%d, state=%s\n", todo.ID, todo.State)

	todo, err = service.UpdateTask(todo.ID, Todo{
		Note: Note{
			Title: "Complete Go assignment",
			Body:  "Implement processing, tests, and documentation",
		},
		Priority:         PriorityHigh,
		EstimatedMinutes: 90,
	})
	if err != nil {
		log.Printf("Could not update task %d: %v", todo.ID, err)
	} else {
		fmt.Printf("Updated to-do: ID=%d, state=%s\n", todo.ID, todo.State)
	}

	todo, err = service.GetTask(todo.ID)
	if err != nil {
		log.Printf("Could not retrieve task %d: %v", todo.ID, err)
	} else {
		fmt.Printf("Retrieved to-do: ID=%d, state=%s\n", todo.ID, todo.State)
	}

	fmt.Println("\nProcessing tasks...")

	for _, task := range store.List() {
		processed, err := service.ProcessTask(task.ID)
		if err != nil {
			log.Printf("Could not process task %d: %v", task.ID, err)
			continue
		}

		fmt.Printf(
			"ID=%d, type=%s, state=%s, result=%+v\n",
			processed.ID,
			processed.Type,
			processed.State,
			processed.Result,
		)
	}

	fmt.Println("\nListing all tasks:")

	allTasks, err := service.GetAllTasks()
	if err != nil {
		log.Printf("Could not list tasks: %v", err)
		return
	}

	for _, task := range allTasks {
		fmt.Printf(
			"ID=%d, type=%s, payload=%+v, state=%s\n",
			task.ID,
			task.Type,
			task.Payload,
			task.State,
		)
	}

	fmt.Println("Total tasks:", len(allTasks))

}
