package auth

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func signinPage(c *fiber.Ctx) error {
	page := signinView()
	handler := adaptor.HTTPHandler(templ.Handler(page))

	return handler(c)
}

func signupPage(c *fiber.Ctx) error {
	page := signupView()
	handler := adaptor.HTTPHandler(templ.Handler(page))

	return handler(c)
}
