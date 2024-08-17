package home

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func homePage(c *fiber.Ctx) error {
	page := homeView()
	handler := adaptor.HTTPHandler(templ.Handler(page))

	return handler(c)
}
