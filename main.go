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

	// Public routes
	e.POST("/signup", controller.Signup)
	e.POST("/login", controller.Login)
	e.POST("/upload", controller.Upload)

	e.POST("/uploadlargefile", controller.UploadLargeFile)

	// Protected routes (require JWT authentication)
	e.POST("/my/processes", controller.Processes, middleware.JWT([]byte("secret")))
	e.POST("/my/statistics", controller.Statistics, middleware.JWT([]byte("secret")))

	// Admin routes (require admin authentication)
	e.POST("/Admin/process_by_username", controller.Process_by_username, middleware.JWT([]byte("secret")))
	e.GET("/Admin/statistics", controller.Admin_statistics, middleware.JWT([]byte("secret")))

	// Token refresh endpoint
	e.POST("/refreshtoken", controller.RefreshToken)

	e.POST("/admin/login", controller.Login)

	// Database setup
	// time.Sleep(20)
	db, err := models.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	controller.SetDB(db)

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE}, // Specify the allowed HTTP methods

		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Start the server
	e.Logger.Fatal(e.Start(":1303"))
}
