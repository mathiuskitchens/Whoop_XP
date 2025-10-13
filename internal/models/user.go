package models

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

func CreateUser(db *sql.DB, username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert new user
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hashed))
	return err
}

func AuthenticateUser(db *sql.DB, username, password string) (int, error) {
	var id int
	var hash string

	err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("invalid username or password")
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, errors.New("invalid username or password")
	}
	return id, nil

}
