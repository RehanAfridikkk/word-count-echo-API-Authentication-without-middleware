package controller

import (
	"net/http"
	"time"

	"github.com/RehanAfridikkk/API-Authentication/structure"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "Rehan" || password != "Bahi!" {
		return echo.ErrUnauthorized
	}

	claims := &structure.JwtCustomClaims{
		Name:  "Rehan Afrdi",
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30000)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
