package dashboard

import (
	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func dashboardPage(c *fiber.Ctx) error {
	user, err := database.GetUserFromCookie(c)
	rootId := c.Params("rootId")

	if rootId == "" {
		rootId = "root"
	} else {
		if err != nil {
			c.Set("HX-Redirect", "/auth/signin")
			return c.Status(fiber.StatusInternalServerError).Redirect("/auth/signin")
		}

		// Check if the document exists and if it's a directory
		doc, err := database.GetDocumentFromId(rootId)
		if err != nil {
			c.Set("HX-Redirect", "/dashboard/root")
			return c.Status(fiber.StatusNotFound).Redirect("/dashboard")
		}

		if doc.Type != "directory" {
			c.Set("HX-Redirect", "/dashboard/root")
			return c.Status(fiber.StatusNotFound).Redirect("/dashboard")
		}
	}

	page := dashboardView(user, rootId)
	handler := adaptor.HTTPHandler(templ.Handler(page))

	return handler(c)
}
