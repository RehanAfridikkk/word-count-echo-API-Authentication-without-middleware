package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UploadRequest struct {
	ID      string `json:"ID"`
	RunTime string `json:"runtime"`
}

func Processes(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token is missing"})
	}

	tok := strings.Split(authHeader, " ")
	if len(tok) != 2 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
	}

	tokenString := tok[1]

	username, err := Validate(tokenString)
	if err != nil {
		return handleTokenError(c, err)
	}
	jsonResult := QueryProcesses(username)

	return c.JSON(http.StatusOK, echo.Map{
		"result": jsonResult,
	})

}

func QueryProcesses(username string) interface{} {
	// Create a slice to store multiple instances of UploadRequest
	var uploadRequests []UploadRequest

	fmt.Println(username)
	rows, err := db.Query(`
        SELECT id, run_time FROM upload_requests WHERE username = $1
    `, username)
	if err != nil {
		// Handle the error
		log.Fatal(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		// Create a new instance of the struct for each row
		var uploadRequest UploadRequest

		if err := rows.Scan(&uploadRequest.ID, &uploadRequest.RunTime); err != nil {
			// Handle the scan error
			log.Fatal(err)
			return nil
		}

		// Append the struct to the slice
		uploadRequests = append(uploadRequests, uploadRequest)
	}

	if err := rows.Err(); err != nil {
		// Handle the error
		log.Fatal(err)
		return nil
	}

	// Return the slice directly without encoding to JSON
	return uploadRequests
}
