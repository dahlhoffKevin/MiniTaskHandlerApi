package main

import (
	"fmt"
	"log"
	"net/http"

	"go-task-api/handlers"
	"go-task-api/types"
	"go-task-api/utils"
)

func main() {
	store := &types.InMemoryTaskStore{
		Tasks:  []types.Task{},
		NextID: 1,
	}

	taskHandler := handlers.NewTaskHandler(store)

	// register routes
	http.HandleFunc("/", handleHelloWorld)

	// handle tasks
	http.HandleFunc("/tasks", utils.RouteLogging(taskHandler.HandleTasks))
	http.HandleFunc("/tasks/{id}", utils.RouteLogging(taskHandler.HandleTasks))
	http.HandleFunc("/tasks/{id}/done", utils.RouteLogging(taskHandler.HandleTasks))
	http.HandleFunc("/tasks/{id}/rename", utils.RouteLogging(taskHandler.HandleTasks))

	fmt.Println("Server läuft auf :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
	utils.LogToConsole("hanling /")
	w.WriteHeader(http.StatusOK)
}
