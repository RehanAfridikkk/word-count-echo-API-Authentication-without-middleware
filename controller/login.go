// controller/user.go

package controller

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RehanAfridikkk/API-Authentication/models"
	"github.com/RehanAfridikkk/API-Authentication/structure"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

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

	// Generate a refresh token
	refreshClaims := &structure.JwtCustomClaims{
		Name:  username,
		Admin: role == "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Refresh token expires in 7 days
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	rt, err := refreshToken.SignedString([]byte("refresh_secret"))
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating refresh token")
	}

	claims := &structure.JwtCustomClaims{
		Name:  username,
		Admin: role == "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
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
		"token":         t,
		"refresh_token": rt,
	})
}

// RefreshToken handles the token refresh request.
func RefreshToken(c echo.Context) error {
	refreshToken := c.Request().Header.Get("Authorization")
	if refreshToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Refresh token is required")
	}

	// Extract the token from the "Bearer" prefix
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	// Validate the refresh token
	refreshClaims := new(structure.JwtCustomClaims)
	refreshTokenSecret := []byte("refresh_secret") // Replace with secure secret management
	rt, err := jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenSecret, nil
	})
	if err != nil {
		log.Println("Error parsing refresh token:", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}
	if !rt.Valid {
		log.Println("Refresh token is not valid")
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	// Generate a new access token
	accessTokenSecret := []byte("secret") // Replace with secure secret management
	claims := &structure.JwtCustomClaims{
		Name:  refreshClaims.Name,
		Admin: refreshClaims.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // New token expires in 24 hours
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := newToken.SignedString(accessTokenSecret)
	if err != nil {
		log.Println("Error generating new token:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating new token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// comparePasswords compares the provided password with the hashed password.
func comparePasswords(providedPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		log.Println("Password comparison error:", err)
	}
	return err == nil
}

// ... (rest of the code)
