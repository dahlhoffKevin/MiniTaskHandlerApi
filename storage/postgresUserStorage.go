package storage

import (
	"database/sql"
	"go-task-api/httpError"
	"go-task-api/types"
	"net/http"

	"github.com/google/uuid"
)

type PostgresqlUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresqlUserStore {
	return &PostgresqlUserStore{db: db}
}

func (s *PostgresqlUserStore) GetAll() ([]types.User, *httpError.HTTPError) {
	rows, err := s.db.Query(`SELECT id, name, email FROM users`)
	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "failed to query user: "+err.Error())
	}
	defer rows.Close()

	users := make([]types.User, 0)
	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, httpError.New(http.StatusInternalServerError, "faild to scan user: "+err.Error())
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "row iteration failed"+err.Error())
	}

	if len(users) == 0 {
		return nil, httpError.New(http.StatusNotFound, "no users found")
	}

	return users, nil
}

func (s *PostgresqlUserStore) GetByID(id uuid.UUID) (*types.User, *httpError.HTTPError) {
	var u types.User

	err := s.db.QueryRow(`SELECT id, name, email FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Name, &u.Email)

	if err == sql.ErrNoRows {
		return nil, httpError.New(http.StatusNotFound, "user not found")
	}

	if err != nil {
		return nil, httpError.New(http.StatusInternalServerError, "failed to query user: "+err.Error())
	}

	return &u, nil
}

func (s *PostgresqlUserStore) Create(name string, email string) (types.User, *httpError.HTTPError) {
	var u types.User

	err := s.db.QueryRow(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email`, name, email).
		Scan(&u.ID, &u.Name, &u.Email)

	if err != nil {
		return types.User{}, httpError.New(http.StatusInternalServerError, "failed to create new user: "+err.Error())
	}

	return u, nil
}

func (s *PostgresqlUserStore) Delete(id uuid.UUID) *httpError.HTTPError {
	res, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return httpError.New(http.StatusInternalServerError, "failed to delete user: "+err.Error())
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return httpError.New(http.StatusInternalServerError, "failed to read rows affected: "+err.Error())
	}
	if affected == 0 {
		return httpError.New(http.StatusNotFound, "user not found")
	}

	return nil
}
