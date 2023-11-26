package main

import (
	"log"

	"github.com/RehanAfridikkk/API-Authentication/controller"
	"github.com/RehanAfridikkk/API-Authentication/models"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/login", controller.Login)
	e.GET("/", controller.Accessible)
	e.POST("/upload", controller.Upload)
	e.POST("/signup", controller.Signup)

	db, err := models.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	controller.SetDB(db)

	models.PingDB(db)

	e.Logger.Fatal(e.Start(":1303"))

}
