// controller/user.go

package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/RehanAfridikkk/API-Authentication/models"
	"github.com/RehanAfridikkk/API-Authentication/structure"
	"github.com/labstack/echo/v4"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

// Signup handles user signup.
func Signup(c echo.Context) error {
	// Parse the request body to get user information
	user := new(structure.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Set the default role for users
	if user.Role == "" {
		user.Role = "user"
	}

	// Validate input
	if err := validateUser(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if the user already exists based on user ID
	existingUser, err := models.GetUserByusername(db, user.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking user existence")
	}
	if existingUser != nil {
		return echo.NewHTTPError(http.StatusConflict, "User with the provided ID already exists")
	}

	// Create the user
	if err := models.CreateUser(db, user); err != nil {
		log.Println("Error creating user:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	// Additional logic for admin signup
	if user.Role == "admin" {
		isAdminAllowed := checkAdminPermission(db, user.Username)
		if !isAdminAllowed {
			return echo.NewHTTPError(http.StatusUnauthorized, "Admin signup not allowed")
		}

		adminSignupTasks(db, user.Username)
	}

	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

func validateUser(user *structure.User) error {

	if user.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username cannot be empty")
	}
	if len(user.Password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters")
	}

	return nil
}

func adminSignupTasks(db *sql.DB, username string) {

	_, err := db.Exec("UPDATE users SET is_admin = true WHERE username = $1", username)
	if err != nil {
		log.Println("Error granting admin privileges:", err)
		return
	}

	log.Println("Admin privileges granted for user:", username)
}

func checkAdminPermission(db *sql.DB, username string) bool {

	var role string
	err := db.QueryRow("SELECT role FROM users WHERE id = $1", username).Scan(&role)
	if err == sql.ErrNoRows {

		log.Println("User not found")
		return false
	} else if err != nil {
		log.Println("Error checking admin permission:", err)
		return false
	}

	return role == "admin"
}
