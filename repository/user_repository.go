package repository

import (
	"database/sql"
	"errors"
)

// User represents a user in the system
type User struct {
	ID       int
	Username string
	Password string
}

// UserRepository defines methods for interacting with the users data.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository initializes and returns a new UserRepositorys
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetUserByUsername fetches a user from the database by username.
func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := r.DB.QueryRow(query, username)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}
