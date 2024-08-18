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
	}

	if err != nil {
		c.Set("HX-Redirect", "/auth/signin")
		return c.Status(fiber.StatusInternalServerError).Redirect("/auth/signin")
	}

	page := dashboardView(user, rootId)
	handler := adaptor.HTTPHandler(templ.Handler(page))

	return handler(c)
}
