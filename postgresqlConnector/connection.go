package postgresqlConnector

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"go-task-api/httpError"

	_ "github.com/lib/pq"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateInitialDatabaseConnection() (*sql.DB, error) {
	connStr := "postgres://dev:devpass@localhost:5432/gotaskapi?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "failed to open database connection: "+err.Error())
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, httpError.New(http.StatusInternalServerError, "failed to ping database: "+err.Error())
	}

	log.Println("database connection successful")
	return db, nil
}

func TestConnection(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.QueryContext(r.Context(), `SELECT id, name, email FROM users`)
	if err != nil {
		httpError.Write(w, httpError.New(http.StatusInternalServerError, "failed to query database: "+err.Error()))
		return
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			httpError.Write(w, httpError.New(http.StatusInternalServerError, "failed to scan database row: "+err.Error()))
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		httpError.Write(w, httpError.New(http.StatusInternalServerError, "row iteration failed: "+err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
