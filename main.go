package main

import (
	"log"

	"github.com/RehanAfridikkk/API-Authentication/controller"
	"github.com/RehanAfridikkk/API-Authentication/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.POST("/login", controller.Login)
	e.POST("/upload", controller.Upload)
	e.POST("/signup", controller.Signup)
	e.POST("/my/processes", controller.Processes)
	e.GET("/my/statistics", controller.Statistics)
	e.POST("/Admin/process_by_username", controller.Process_by_username)
	e.GET("/Admin/statistics", controller.Admin_statistics)
	e.POST("/refreshtoken", controller.RefreshToken)

	db, err := models.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	controller.SetDB(db)

	models.PingDB(db)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Logger.Fatal(e.Start(":1303"))

}
