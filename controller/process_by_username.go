package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func Process_by_username(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token is missing"})
	}

	tok := strings.Split(authHeader, " ")
	if len(tok) != 2 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
	}

	tokenString := tok[1]

	claims, err := ValidateAdmin(tokenString)
	if err != nil {
		return handleTokenError(c, err)
	}

	// Check if the user has admin role
	if !claims {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Only admin can check the data"})
	}

	// Parse the username from the request body
	var requestBody struct {
		Username string `json:"username"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	username := requestBody.Username

	// Fetch data for the specified username
	jsonResult := QueryProcesses(username)

	return c.JSON(http.StatusOK, echo.Map{
		"result": jsonResult,
	})
}
func ValidateAdmin(tokenString string) (bool, error) {
	hmacSampleSecret := []byte("secret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		admin, ok := claims["admin"].(bool)
		if !ok {
			return false, fmt.Errorf("admin claim is not a boolean")
		}
		return admin, nil
	} else {
		return false, err
	}
}
