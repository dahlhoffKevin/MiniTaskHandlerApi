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
	// routing setup
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.RouteLogging(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	setupTaskRestEndpoints(mux)
	setupUserRestEndpoints(mux)

	fmt.Println("Server läuft auf :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func setupUserRestEndpoints(mux *http.ServeMux) {
	// storage setup
	store := &storage.InMemoryUserStore{
		Users:  []types.User{},
		NextID: 1,
	}

	userHandler := handlers.NewUserStore(store)

	// handle users
	mux.HandleFunc("/users", utils.RouteLogging(userHandler.HandleUsers))
	mux.HandleFunc("/users/{id}", utils.RouteLogging(userHandler.HandleUsers))
}

func setupTaskRestEndpoints(mux *http.ServeMux) {
	// storage setup
	store := &storage.InMemoryTaskStore{
		Tasks:  []types.Task{},
		NextID: 1,
	}

	taskHandler := handlers.NewTaskHandler(store)

	// handle tasks
	mux.HandleFunc("/tasks", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}/done", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}/rename", utils.RouteLogging(taskHandler.HandleTasks))
}
