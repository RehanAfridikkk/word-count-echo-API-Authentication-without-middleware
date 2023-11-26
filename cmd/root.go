package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/RehanAfridikkk/word-count-Echo-API-fileupload/pkg"
)

func ProcessFile(file multipart.File, routines int) (pkg.CountsResult, int, time.Duration, error) {
	start := time.Now()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return pkg.CountsResult{}, 0, 0, err
	}

	fileContent := buf.Bytes()

	chunkSize := len(fileContent) / routines
	results := make(chan pkg.CountsResult, routines)

	for i := 0; i < routines; i++ {
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == routines-1 {
			endIndex = len(fileContent)
		}

		chunk := fileContent[startIndex:endIndex]
		go pkg.Counts(chunk, results)
	}

	totalCounts := pkg.CountsResult{}

	for i := 0; i < routines; i++ {
		result := <-results
		totalCounts.LineCount += result.LineCount
		totalCounts.WordsCount += result.WordsCount
		totalCounts.VowelsCount += result.VowelsCount
		totalCounts.PunctuationCount += result.PunctuationCount
	}

	runTime := time.Since(start)

	fmt.Println("Number of lines:", totalCounts.LineCount)
	fmt.Println("Number of words:", totalCounts.WordsCount)
	fmt.Println("Number of vowels:", totalCounts.VowelsCount)
	fmt.Println("Number of punctuation:", totalCounts.PunctuationCount)
	fmt.Println("Run Time:", runTime)

	return totalCounts, routines, runTime, nil
}
