package main

import (
	"fmt"

	_ "github.com/Pyramakerz/Library_Management_System/docs"

	config "github.com/Pyramakerz/Library_Management_System/PKG/Config"
	models "github.com/Pyramakerz/Library_Management_System/PKG/Models"
	routes "github.com/Pyramakerz/Library_Management_System/PKG/Routes"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Library Management System
// @version 1.0
// @description This API System is for Pyramakerz company having a CRUD operaions for Author and Book

// @host localhost:9090
// @BasePath /
func main() {
	// 1) Start the Db Connection and Auto Migrate the models to be tables
	config.Connect()
	db := config.GetDB()

	if db == nil {
		fmt.Printf("Failed to connect to the database.")
	}

	err := db.AutoMigrate(&models.Book{}, &models.Author{})
	if err != nil {
		fmt.Printf("Failed to migrate models: %v", err)
	}

	// 2) Set the routes
	app := fiber.New()

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	routes.Library_Management_System_Routes(app)
	app.Listen("localhost:9090")
}
