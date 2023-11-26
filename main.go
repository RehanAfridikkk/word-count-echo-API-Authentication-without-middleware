package main

import (
	"github.com/RehanAfridikkk/API-Authentication/controller"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/login", controller.Login)
	e.GET("/", controller.Accessible)
	e.POST("/upload", controller.Upload)

	e.Logger.Fatal(e.Start(":1303"))
}
