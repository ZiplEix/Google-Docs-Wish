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
	authGroup.Post("/signout", signout)
	// authGroup.Get("/google", googleLogin)
	// authGroup.Get("/google/callback", googleCallback)
}

func signout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "auth_token",
		Value: "",
	})
	c.Set("HX-Redirect", "/")
	return c.SendStatus(fiber.StatusOK)
}
