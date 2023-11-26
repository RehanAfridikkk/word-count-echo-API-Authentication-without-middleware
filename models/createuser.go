// models/user.go

package models

import (
	"database/sql"
	"log"

	"github.com/RehanAfridikkk/API-Authentication/structure"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user model.

// CreateUser creates a new user.
func CreateUser(db *sql.DB, user *structure.User) error {
	// Hash the user's password before storing it in the database
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hashedPassword)
	return err
}

// hashPassword hashes the password using bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
