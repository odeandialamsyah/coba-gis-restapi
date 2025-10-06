package routes

import (
	"pretest-gis/controllers"

	"github.com/gofiber/fiber/v2"
)

func PlaceRoutes(app *fiber.App) {
	api := app.Group("/api")
	places := api.Group("/places")

	places.Post("/", controllers.CreatePlace)
	places.Get("/", controllers.GetPlaces)
	places.Get("/:id", controllers.GetPlace)
	places.Put("/:id", controllers.UpdatePlace)
	places.Delete("/:id", controllers.DeletePlace)
}
