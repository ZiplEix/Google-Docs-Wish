package dashboard

import (
	"github.com/ZiplEix/Google-Docs-Wish/middleware"
	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/dashboard", middleware.Protected())
	dashboardGroup.Get("/", dashboardPage)
}
