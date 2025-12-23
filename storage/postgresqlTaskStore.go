package storage

import (
	"database/sql"
	"go-task-api/httpError"
	"go-task-api/types"
	"net/http"

	"github.com/google/uuid"
)

type PostgresqlTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresqlTaskStore {
	return &PostgresqlTaskStore{db}
}

func (s *PostgresqlTaskStore) GetAll() ([]types.Task, *httpError.HTTPError) {
	rows, err := s.db.Query(`SELECT id, title, done, userid FROM tasks`)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "faild to query tasks: "+err.Error())
	}
	defer rows.Close()

	tasks := make([]types.Task, 0)
	for rows.Next() {
		var t types.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.UserID); err != nil {
			return nil, httpError.New(http.StatusInternalServerError, "failed to scan task: "+err.Error())
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "row iteration failed: "+err.Error())
	}

	if len(tasks) == 0 {
		return nil, httpError.New(http.StatusNotFound, "no tasks found")
	}

	return tasks, nil
}

func (s *PostgresqlTaskStore) GetByID(id uuid.UUID) (*types.Task, *httpError.HTTPError) {
	var t types.Task

	err := s.db.QueryRow(`SELECT id, title, done, userid FROM tasks WHERE id = $1`, id).
		Scan(&t.ID, &t.Title, &t.Done, &t.UserID)

	if err == sql.ErrNoRows {
		return nil, httpError.New(http.StatusNotFound, "task not found")
	}

	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "failed to query task: "+err.Error())
	}

	return &t, nil
}

func (s *PostgresqlTaskStore) Create(title string, userid uuid.UUID) (types.Task, *httpError.HTTPError) {
	var t types.Task
	var u types.User

	//check if user exists
	userCheckError := s.db.QueryRow(`SELECT id, name, email FROM users WHERE id = $1`, userid).
		Scan(&u.ID, &u.Name, &u.Email)

	if userCheckError == sql.ErrNoRows || u.ID != userid {
		return types.Task{}, httpError.New(http.StatusNotFound, "user does not exist")
	}

	err := s.db.QueryRow(`INSERT INTO tasks (title, done, userid) VALUES ($1, $2, $3) RETURNING id, title, done, userid`, title, false, userid).
		Scan(&t.ID, &t.Title, &t.Done, &t.UserID)

	if err != nil {
		return types.Task{}, httpError.New(http.StatusInternalServerError, "failed to create new task: "+err.Error())
	}

	return t, nil
}

func (s *PostgresqlTaskStore) Delete(id uuid.UUID) *httpError.HTTPError {
	res, err := s.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, "failed to delete task: "+err.Error())
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return httpError.New(http.StatusInternalServerError, "failed to read rows affected: "+err.Error())
	}
	if affected == 0 {
		return httpError.New(http.StatusNotFound, "task not found")
	}

	return nil
}

func (s *PostgresqlTaskStore) Update(task types.Task) *httpError.HTTPError {
	var t types.Task

	err := s.db.QueryRow(
		`UPDATE tasks 
		 SET title = $1, done = $2, userid = $3 
		 WHERE id = $4 
		 RETURNING id, title, done, userid`,
		task.Title, task.Done, task.UserID, task.ID,
	).Scan(&t.ID, &t.Title, &t.Done, &t.UserID)

	if err == sql.ErrNoRows {
		return httpError.New(http.StatusNotFound, "task not found")
	}
	if err != nil {
		return httpError.New(http.StatusInternalServerError, "failed to update task: "+err.Error())
	}

	return nil
}
