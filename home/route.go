package home

import (
	"github.com/gofiber/fiber/v2"
)

func HomeRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/")
	dashboardGroup.Get("/", homePage)
}
