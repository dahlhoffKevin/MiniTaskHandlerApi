package main

import (
	"fmt"
	"log"
	"net/http"

	"go-task-api/handlers"
	"go-task-api/storage"
	"go-task-api/types"
	"go-task-api/utils"
)

func main() {
	// storage setup
	store := &storage.InMemoryTaskStore{
		Tasks:  []types.Task{},
		NextID: 1,
	}

	taskHandler := handlers.NewTaskHandler(store)

	// routing setup
	mux := http.NewServeMux()

	mux.HandleFunc("/", utils.RouteLogging(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// handle tasks
	mux.HandleFunc("/tasks", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}/done", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}/rename", utils.RouteLogging(taskHandler.HandleTasks))

	fmt.Println("Server läuft auf :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
