package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RehanAfridikkk/API-Authentication/cmd"
	"github.com/RehanAfridikkk/API-Authentication/structure"

	// "github.com/RehanAfridikkk/word-count-Echo-API/pkg"
	// "github.com/RehanAfridikkk/word-count-Echo-API/pkg"
	// "github.com/RehanAfridikkk/word-count-Echo-API/pkg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Upload(c echo.Context) error {

	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token is missing")
	}

	tok := strings.Split(authHeader, " ")

	fmt.Println("TOKEN: ", authHeader)
	tokenString := tok[1]

	if tokenString == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	fmt.Println("ERROR: ", err)
	if err != nil {
		return err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "File not provided",
		})
	}

	routines, err := strconv.Atoi(c.FormValue("routines"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid value for routines",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to open file",
		})
	}
	defer fileContent.Close()

	totalCounts, _, runTime, err := cmd.ProcessFile(fileContent, routines)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to process file",
		})
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		username := claims["name"].(string)
		fmt.Println(username)

		// Store the upload request information in the database
		err := storeUploadRequest(username, file.Filename, routines, totalCounts, runTime)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to store upload request information",
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"lineCount":        totalCounts.LineCount,
		"wordsCount":       totalCounts.WordsCount,
		"vowelsCount":      totalCounts.VowelsCount,
		"punctuationCount": totalCounts.PunctuationCount,
		"runTime":          runTime.String(),
	})
}

func storeUploadRequest(username, fileName string, routines int, totalCounts structure.CountsResult, runTime time.Duration) error {
	_, err := db.Exec(`

        INSERT INTO upload_requests (username, file_name, routines, line_count, words_count, vowels_count, punctuation_count, run_time)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, username, fileName, routines, totalCounts.LineCount, totalCounts.WordsCount, totalCounts.VowelsCount, totalCounts.PunctuationCount, runTime)
	return err
}
