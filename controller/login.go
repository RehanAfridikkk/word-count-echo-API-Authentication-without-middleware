// controller/user.go

package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/RehanAfridikkk/API-Authentication/models"
	"github.com/RehanAfridikkk/API-Authentication/structure"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login.
// controller/user.go

// ...

// Login handles user login.
func Login(c echo.Context) error {
	loginRequest := new(structure.LoginRequest)
	if err := c.Bind(loginRequest); err != nil {
		log.Println("Error binding login request:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	username := loginRequest.Username
	password := loginRequest.Password

	log.Println("Login request:", loginRequest)

	// Validate the username and password against the database
	user, err := models.GetUserByusername(db, loginRequest.Username)
	if err != nil {
		log.Println("Error checking credentials:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking credentials")
	}

	if user == nil {
		log.Println("User not found for ID:", loginRequest.Username)
		return echo.ErrUnauthorized
	}

	log.Println("Found user:", user)

	// Compare the provided password with the hashed password stored in the database
	if !comparePasswords(password, user.Password) {
		log.Println("Password mismatch for user:", username)
		return echo.ErrUnauthorized
	}

	// Assume the role is "admin" for the admin user
	role := "admin"
	if user.Role != "admin" {
		role = "user"
	}

	claims := &structure.JwtCustomClaims{
		Name:  username,
		Admin: role == "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println("Error generating token:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	log.Println("Login successful for user:", username)

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// ...

// comparePasswords compares the provided password with the hashed password.
func comparePasswords(providedPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		log.Println("Password comparison error:", err)
	}
	return err == nil
}

// ... rest of the code ...
