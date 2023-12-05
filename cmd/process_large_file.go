package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/RehanAfridikkk/API-Authentication/pkg"
	"github.com/RehanAfridikkk/API-Authentication/structure"
	"github.com/gorilla/websocket"
)

// WebSocketHandler is the handler for WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}
	defer connection.Close()

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	routinesStr := r.FormValue("routines")
	routines, err := strconv.Atoi(routinesStr)
	if err != nil {
		http.Error(w, "Could not get routines from form", http.StatusBadRequest)
		return
	}

	totalCounts, processedRoutines, runTime, err := ProcessLargeFile(file, routines, connection)
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}

	fmt.Println("Total Counts:", totalCounts)
	fmt.Println("Number of Routines:", processedRoutines)
	fmt.Println("Run Time:", runTime)
}

// ProgressUpdate represents a progress update message
type ProgressUpdate struct {
	Percent int    `json:"percent"`
	Message string `json:"message"`
}

// CompletionMessage represents a completion message
type CompletionMessage struct {
	Message string `json:"message"`
}

// ProcessFile processes a file in chunks using goroutines and sends progress updates via WebSocket
func ProcessLargeFile(file multipart.File, routines int, conn *websocket.Conn) (structure.CountsResult, int, time.Duration, error) {
	start := time.Now()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return structure.CountsResult{}, 0, 0, err
	}

	fileContent := buf.Bytes()

	chunkSize := len(fileContent) / routines
	results := make(chan structure.CountsResult, routines)

	for i := 0; i < routines; i++ {
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == routines-1 {
			endIndex = len(fileContent)
		}

		chunk := fileContent[startIndex:endIndex]

		// Send progress update to the client
		progressUpdate := ProgressUpdate{
			Percent: (i + 1) * 100 / routines,
			Message: fmt.Sprintf("Processing chunk %d", i+1),
		}
		conn.WriteJSON(progressUpdate)

		go pkg.Counts(chunk, results)
	}

	totalCounts := structure.CountsResult{}

	for i := 0; i < routines; i++ {
		result := <-results
		totalCounts.LineCount += result.LineCount
		totalCounts.WordsCount += result.WordsCount
		totalCounts.VowelsCount += result.VowelsCount
		totalCounts.PunctuationCount += result.PunctuationCount
	}

	runTime := time.Since(start)

	// Send completion message to the client
	completionMessage := CompletionMessage{
		Message: "File processing completed!",
	}
	conn.WriteJSON(completionMessage)

	return totalCounts, routines, runTime, nil
}
