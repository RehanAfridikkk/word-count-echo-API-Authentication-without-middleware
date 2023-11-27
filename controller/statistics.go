package controller

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type StatisticsResult struct {
	ExecutionCount big.Int       `json:"execution_count"`
	AverageRuntime time.Duration `json:"average_runtime"`
}

// Statistics handles the statistics endpoint.
func Statistics(c echo.Context) error {
	// Extract the token string from the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token is missing"})
	}

	// Split the token string to get the actual token
	tok := strings.Split(authHeader, " ")
	if len(tok) != 2 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
	}

	tokenString := tok[1]

	// Validate the token and extract user information
	username, err := Validate(tokenString)
	if err != nil {
		return handleTokenError(c, err)
	}

	// Get file name from the request
	fileName := c.QueryParam("file")
	if fileName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "File name is required",
		})
	}

	// Query the database for statistics
	stats, err := queryStatistics(username, fileName)
	if err != nil {
		return handleStatisticsError(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"username":        username,
		"file_name":       fileName,
		"execution_count": stats.ExecutionCount.String(), // Convert big.Int to string
		"average_runtime": stats.AverageRuntime.String(),
	})
}

func handleTokenError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *echo.HTTPError:
		return c.JSON(e.Code, echo.Map{"error": e.Message})
	default:
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to extract username from token"})
	}
}

func handleStatisticsError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *echo.HTTPError:
		return c.JSON(e.Code, echo.Map{"error": e.Message})
	default:
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve statistics"})
	}
}

func Validate(tokenString string) (string, error) {
	hmacSampleSecret := []byte("secret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println("line87")
		return hmacSampleSecret, nil
	})
	fmt.Println("line90")
	if err != nil {
		fmt.Println("line92")
		return "", err

	}
	fmt.Println("line97")
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		user, ok := claims["name"].(string)
		if !ok {
			fmt.Println("line100")
			return "", fmt.Errorf("user claim is not a float64")
		}
		fmt.Println("line103")
		return user, nil
	} else {
		fmt.Println("line101")
		return "", err
	}
}

// queryStatistics queries the database for statistics.
func queryStatistics(username, fileName string) (StatisticsResult, error) {
	// Perform a database query to retrieve execution count and average runtime
	// for the given username and file name from the upload_requests table

	// Replace the following with your actual database query
	fmt.Println(username, fileName)
	row := db.QueryRow(`
	SELECT COUNT(*) AS execution_count, AVG(run_time) AS average_runtime FROM upload_requests WHERE username = $1 AND file_name = $2
	`, username, fileName)

	var executionCount big.Int
	var averageRuntimeSeconds float64

	err := row.Scan(&executionCount, &averageRuntimeSeconds)
	if err != nil {
		// Handle the error (e.g., user not found, database error)
		return StatisticsResult{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve statistics from the database")
	}

	averageRuntime := time.Duration(int64(averageRuntimeSeconds)) * time.Second

	return StatisticsResult{
		ExecutionCount: executionCount,
		AverageRuntime: averageRuntime,
	}, nil
}
