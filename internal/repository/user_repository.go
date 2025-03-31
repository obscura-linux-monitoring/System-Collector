package repository

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsUserByObscuraKey(ObscuraKey string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE obscura_key = $1)`
	row := r.db.QueryRow(query, ObscuraKey)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
