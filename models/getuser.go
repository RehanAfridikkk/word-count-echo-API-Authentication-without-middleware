// models/user.go

package models

import (
	"database/sql"
	"log"

	"github.com/RehanAfridikkk/API-Authentication/structure"
)

// GetUserByID retrieves a user by user ID.
func GetUserByID(db *sql.DB, userID int) (*structure.User, error) {
	var user structure.User
	err := db.QueryRow("SELECT id, username, password, role FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		// User not found
		return nil, nil
	} else if err != nil {
		log.Println("Error querying user by ID:", err)
		return nil, err
	}

	return &user, nil
}
