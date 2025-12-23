package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go-task-api/handlers"
	"go-task-api/postgresqlConnector"
	"go-task-api/storage"
	"go-task-api/utils"
)

func main() {
	// routing setup
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.RouteLogging(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	db, err := postgresqlConnector.CreateInitialDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux.HandleFunc("/sqlTest", utils.RouteLogging(func(w http.ResponseWriter, r *http.Request) {
		postgresqlConnector.TestConnection(db, w, r)
	}))

	setupTaskRestEndpoints(mux, db)
	setupUserRestEndpoints(mux, db)

	fmt.Println("Server läuft auf :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func setupUserRestEndpoints(mux *http.ServeMux, db *sql.DB) {
	// storage setup
	store := storage.NewPostgresUserStore(db)

	userHandler := handlers.NewUserStore(store)

	// handle users
	mux.HandleFunc("/users", utils.RouteLogging(userHandler.HandleUsers))
	mux.HandleFunc("/users/{id}", utils.RouteLogging(userHandler.HandleUsers))
}

func setupTaskRestEndpoints(mux *http.ServeMux, db *sql.DB) {
	// storage setup
	store := storage.NewPostgresTaskStore(db)

	taskHandler := handlers.NewTaskHandler(store)

	// handle tasks
	mux.HandleFunc("/tasks", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}/done", utils.RouteLogging(taskHandler.HandleTasks))
	mux.HandleFunc("/tasks/{id}/rename", utils.RouteLogging(taskHandler.HandleTasks))
}
