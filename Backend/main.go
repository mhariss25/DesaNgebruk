package main

import (
	"github.com/agnarbriantama/DesaNgembruk-Backend/database"
	"github.com/agnarbriantama/DesaNgembruk-Backend/database/migrations"
	"github.com/agnarbriantama/DesaNgembruk-Backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Inisialisasi koneksi ke database dan melakukan auto migrate
	database.InitDatabase()
	migrations.MigrationTable()
	// awsutils.CreateBucket()

	app := fiber.New()

	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.SetupRoutes(app)

	err := app.Listen(":8080")
	if err != nil {
		panic("Failed to start server")
	}
}
