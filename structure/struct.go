package structure



import (
	"mime/multipart"
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