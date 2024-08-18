package document

import (
	"fmt"

	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/ZiplEix/Google-Docs-Wish/middleware"
	"github.com/gofiber/fiber/v2"
)

func DocumentRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/document", middleware.Protected())
	dashboardGroup.Get("/", redirectToDashboard)
	dashboardGroup.Get("/:docId", documentPage)
	dashboardGroup.Post("/create-new", createNewDocument)
	dashboardGroup.Delete("/:docId", deleteDocument)
}

func redirectToDashboard(c *fiber.Ctx) error {
	c.Set("HX-Redirect", "/dashboard")
	return c.Redirect("/dashboard")
}

func documentPage(c *fiber.Ctx) error {
	docId := c.Params("docId")

	return c.SendString("Document: " + docId)
}

func createNewDocument(c *fiber.Ctx) error {
	userId := c.Locals("userID").(string)

	doc, err := database.CreateNewDocInDb(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	fmt.Println("Created new document with ID:", doc.ID)

	c.Set("HX-Redirect", "/document/"+doc.ID)
	return c.SendStatus(fiber.StatusCreated)
}

func deleteDocument(c *fiber.Ctx) error {
	docId := c.Params("docId")

	err := database.DeleteDocumentById(docId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).Send(nil)
}
