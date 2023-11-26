// controller/user.go

package controller

import (
	"database/sql"
	"net/http"

	"log"

	"github.com/RehanAfridikkk/API-Authentication/models"
	"github.com/labstack/echo/v4"
)

var db *sql.DB // Declare a variable to hold the database connection

// SetDB sets the db variable from the main package
func SetDB(database *sql.DB) {
	db = database
}

/// Signup handles user signup.
func Signup(c echo.Context) error {
	// Parse the request body to get user information
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
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

	// Do not include the password in the response
	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

// validateUser performs basic validation on the user input
func validateUser(user *models.User) error {
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
