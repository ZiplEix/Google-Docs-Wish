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
	dashboardGroup.Post("/create-new/:rootId", createNewDocument)
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
	rootId := c.Params("rootId")
	tipe := c.FormValue("type")

	doc, err := database.CreateNewDocInDb(userId, rootId, tipe)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	fmt.Println("Created new document with ID:", doc.ID)

	c.Set("HX-Redirect", "/document/"+doc.ID)
	return c.SendStatus(fiber.StatusCreated)
}

func deleteDocument(c *fiber.Ctx) error {
	docId := c.Params("docId")
	userId := c.Locals("userID").(string)

	// get the document to delete
	doc, err := database.GetDocumentFromId(docId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// check if the user is the author of the document
	if doc.UserId != userId {
		return c.Status(fiber.StatusUnauthorized).SendString("You are not the author of this document")
	}

	err = database.DeleteDocumentById(docId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).Send(nil)
}
