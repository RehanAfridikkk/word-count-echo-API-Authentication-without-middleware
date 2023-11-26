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

var db *sql.DB // Declare a variable to hold the database connection

// SetDB sets the db variable from the main package
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
	existingUser, err := models.GetUserByID(db, user.ID)
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
		isAdminAllowed := checkAdminPermission(db, user.ID)
		if !isAdminAllowed {
			// You might want to roll back the signup or handle it accordingly
			return echo.NewHTTPError(http.StatusUnauthorized, "Admin signup not allowed")
		}

		// Perform additional tasks for admin signup
		// For example, grant admin privileges or notify the owner for approval
		adminSignupTasks(db, user.ID)
	}

	// Do not include the password in the response
	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

// validateUser performs basic validation on the user input
func validateUser(user *structure.User) error {
	if user.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	if user.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username cannot be empty")
	}
	if len(user.Password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters")
	}
	// You can add more validation rules based on your requirements
	// For example, check if the username is unique, etc.
	// ...

	return nil
}

// adminSignupTasks performs additional tasks for admin signup.
func adminSignupTasks(db *sql.DB, userID int) {
	// Grant admin-specific privileges
	// For example, update the user record to indicate admin privileges.

	_, err := db.Exec("UPDATE users SET is_admin = true WHERE id = $1", userID)
	if err != nil {
		log.Println("Error granting admin privileges:", err)
		return
	}

	log.Println("Admin privileges granted for user:", userID)
}

// checkAdminPermission checks if a user is allowed to sign up as an admin.
func checkAdminPermission(db *sql.DB, userID int) bool {
	// Implement your logic to check if the user is allowed to sign up as an admin.
	// This could involve querying the database for specific criteria or permissions.

	// For example, check if the user is the owner of the database.
	// You might have a table or a field in the users table that indicates ownership.

	// Placeholder logic:
	var role string
	err := db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
	if err == sql.ErrNoRows {
		// No rows returned, meaning the user does not exist.
		log.Println("User not found")
		return false
	} else if err != nil {
		log.Println("Error checking admin permission:", err)
		return false
	}

	// Check if the user has an "admin" role
	return role == "admin"
}
