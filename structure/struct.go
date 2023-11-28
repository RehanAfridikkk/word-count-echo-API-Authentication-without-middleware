package structure

import (
	"math/big"
	"mime/multipart"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

type Message struct {
	Upload   multipart.File `form:"upload"`
	Routines int            `form:"routines"`
}

type CountsResult struct {
	LineCount        int
	WordsCount       int
	VowelsCount      int
	PunctuationCount int
}

type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type User struct {
	// ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type StatisticsResult struct {
	ExecutionCount big.Int       `json:"execution_count"`
	AverageRuntime time.Duration `json:"average_runtime"`
}

type UploadRequest struct {
	ID               string `json:"ID"`
	RunTime          string `json:"runtime"`
	WordCount        int    `json:"wordCount"`
	VowelsCount      int    `json:"vowelsCount"`
	PunctuationCount int    `json:"punctuationCount"`
	Routines         int    `json:"routines"`
	LineCount        int    `json:"lineCount"`
}

type AdminStatisticsResult struct {
	Username       string       `json:"username"`
	ExecutionCount int64        `json:"execution_count"`
	AverageRuntime JSONDuration `json:"average_runtime"`
	FileName       string       `json:"file_name"`
}
type JSONDuration struct {
	time.Duration
}
