package main

import (
	"log"
	"os"

	"pretest-gis/config"
	"pretest-gis/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()

	app := fiber.New()

	// Routes
	routes.PlaceRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Server running on port " + port)
	log.Fatal(app.Listen(":" + port))
}
