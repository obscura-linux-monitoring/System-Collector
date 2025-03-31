package models

type User struct {
	GoogleID   string `db:"google_id"`
	Email      string `db:"email"`
	Name       string `db:"name"`
	ObscuraKey string `db:"obscura_key"`
}
