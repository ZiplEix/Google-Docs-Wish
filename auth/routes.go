package auth

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")
	authGroup.Get("/signin", signinPage)
	authGroup.Post("/signin", signin)
	authGroup.Get("/signup", signupPage)
	authGroup.Post("/signup", signup)
	// authGroup.Get("/google", googleLogin)
	// authGroup.Get("/google/callback", googleCallback)
}
