package dashboard

import (
	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func getUserFromCookie(c *fiber.Ctx) (database.User, error) {
	userID := c.Locals("userID").(string)

	state, err := database.FirestoreClient.Doc("users/" + userID).Get(c.Context())
	if err != nil {
		return database.User{}, err
	}

	user := database.NewUser(state.Data(), state.Ref.ID)

	return *user, nil
}

func dashboardPage(c *fiber.Ctx) error {
	user, err := getUserFromCookie(c)
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
