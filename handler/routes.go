package handler

import (
	"github.com/ZiplEix/Google-Docs-Wish/auth"
	"github.com/ZiplEix/Google-Docs-Wish/dashboard"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	auth.AuthRoutes(app)
	dashboard.DashboardRoutes(app)
}
