package controller

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Admin_statistics(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token is missing"})
	}

	tok := strings.Split(authHeader, " ")
	if len(tok) != 2 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
	}

	tokenString := tok[1]

	_, err := ValidateAdmin(tokenString)
	if err != nil {
		return handleTokenError(c, err)
	}

	fileName := c.QueryParam("file")
	if fileName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "File name is required",
		})
	}

	// Query the database for statistics
	stats, err := queryStatisticsAdmin(fileName)
	if err != nil {
		return handleStatisticsError(c, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		// "username":        username,
		"file_name":       fileName,
		"execution_count": stats.ExecutionCount.String(), // Convert big.Int to string
		"average_runtime": stats.AverageRuntime.String(),
	})
}

func queryStatisticsAdmin(fileName string) (StatisticsResult, error) {
	fmt.Println("Querying statistics for", fileName)
	row := db.QueryRow(`
        SELECT COUNT(*) AS execution_count, AVG(run_time) AS average_runtime
        FROM upload_requests
        WHERE file_name = $1
    `, fileName)

	var executionCount int64
	var averageRuntimeStr string

	err := row.Scan(&executionCount, &averageRuntimeStr)
	if err != nil {
		fmt.Println("Error scanning row:", err)
		return StatisticsResult{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve statistics from the database")
	}

	fmt.Println("Retrieved average runtime string:", averageRuntimeStr)

	averageRuntime, err := parseDurationAdmin(averageRuntimeStr)
	if err != nil {
		fmt.Println("Error parsing average runtime:", err)
		return StatisticsResult{}, echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse average runtime from the database")
	}

	fmt.Println("Statistics retrieved successfully:", executionCount, averageRuntime)
	return StatisticsResult{
		ExecutionCount: *big.NewInt(executionCount),
		AverageRuntime: JSONDuration{Duration: averageRuntime},
	}, nil
}
func parseDurationAdmin(durationStr string) (time.Duration, error) {
	// Convert the string to a float64
	durationInSeconds, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert the float64 to a time.Duration
	duration := time.Duration(durationInSeconds * float64(time.Second))

	return duration, nil
}
