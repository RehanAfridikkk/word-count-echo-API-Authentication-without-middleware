package controller

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RehanAfridikkk/API-Authentication/structure"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

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
