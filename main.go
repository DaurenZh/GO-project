package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// readLine reads a line of user input, including spaces.
func readLine(text string) string {
	fmt.Print(text)

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// readInt reads an integer from the terminal.
func readInt(text string) (int, error) {
	input := readLine(text)
	return strconv.Atoi(input)
}

// createHabit asks the user for habit information.
func createHabit(service *Service) {
	name := readLine("Habit name: ")
	title := readLine("Title: ")
	body := readLine("Description: ")

	completedDays, err := readInt("Completed days: ")
	if err != nil {
		fmt.Println("Invalid number.")
		return
	}

	totalDays, err := readInt("Total days: ")
	if err != nil {
		fmt.Println("Invalid number.")
		return
	}

	habit := Habit{
		Name: name,
		Note: Note{
			Title: title,
			Body:  body,
		},
		CompletedDays: completedDays,
		TotalDays:     totalDays,
	}

	task, err := service.CreateTask(TypeHabit, habit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Habit created! ID=%d, State=%s, Payload=%+v\n", task.ID, task.State, task.Payload)
}

// createNote asks the user for note information.
func createNote(service *Service) {
	title := readLine("Note title: ")
	body := readLine("Note content: ")

	note := Note{
		Title: title,
		Body:  body,
	}

	task, err := service.CreateTask(TypeNote, note)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Note created! ID=%d, State=%s, Payload=%+v\n", task.ID, task.State, task.Payload)
}

// createTodo asks the user for to-do information.
func createTodo(service *Service) {
	title := readLine("To-do title: ")
	body := readLine("Description: ")
	priority := Priority(strings.ToLower(readLine(
		"Priority (low/medium/high): ",
	)))

	todo := Todo{
		Note: Note{
			Title: title,
			Body:  body,
		},
		Priority:         priority,
	}

	task, err := service.CreateTask(TypeTodo, todo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("To-do created! ID=%d, State=%s, Payload=%+v\n", task.ID, task.State, task.Payload)
}

// listTasks displays all stored tasks.
func listTasks(service *Service) {
	tasks, err := service.GetAllTasks()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("\n=== All Tasks ===")

	for _, task := range tasks {
		fmt.Printf(
			"ID=%d | Type=%s | State=%s | Payload=%+v\n",
			task.ID,
			task.Type,
			task.State,
			task.Payload,
		)
	}
}

// processTask asks which task should be processed.
func processTask(service *Service) {
	id, err := readInt("Enter task ID: ")
	if err != nil {
		fmt.Println("Invalid ID.")
		return
	}

	task, err := service.ProcessTask(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Task processed successfully!")
	fmt.Printf("ID: %d\n", task.ID)
	fmt.Printf("State: %s\n", task.State)
	fmt.Printf("Payload: %+v\n", task.Payload)
	fmt.Printf("Result: %+v\n", task.Result)
}

func main() {
	store := NewStore()
	service := NewService(store)

	for {
		fmt.Println("\n=== LifeTrack ===")
		fmt.Println("1. Create a habit")
		fmt.Println("2. Create a note")
		fmt.Println("3. Create a to-do")
		fmt.Println("4. List all tasks")
		fmt.Println("5. Process a task")
		fmt.Println("0. Exit")

		choice := readLine("Choose an option: ")

		switch choice {
		case "1":
			createHabit(service)

		case "2":
			createNote(service)

		case "3":
			createTodo(service)

		case "4":
			listTasks(service)

		case "5":
			processTask(service)

		case "0":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}