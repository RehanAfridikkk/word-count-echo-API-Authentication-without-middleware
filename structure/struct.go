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
