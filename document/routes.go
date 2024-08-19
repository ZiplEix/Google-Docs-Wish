package document

import (
	"fmt"

	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/ZiplEix/Google-Docs-Wish/middleware"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func DocumentRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/document", middleware.Protected())
	dashboardGroup.Get("/", redirectToDashboard)
	dashboardGroup.Post("/create-new/:rootId", createNewDocument)
	dashboardGroup.Delete("/:docId", deleteDocument)
	dashboardGroup.Get("/rename/rename_modal/:docId", renameModal)
	dashboardGroup.Post("rename/:docId", renameDocument)

	// document edit page
	editDocGroup := dashboardGroup.Group("/e")
	editDocGroup.Get("/:docId", documentPage)
}

func redirectToDashboard(c *fiber.Ctx) error {
	c.Set("HX-Redirect", "/dashboard")
	return c.Redirect("/dashboard")
}

func documentPage(c *fiber.Ctx) error {
	user, err := database.GetUserFromCookie(c)
	if err != nil {
		c.Set("HX-Redirect", "/auth/signin")
		return c.Status(fiber.StatusInternalServerError).Redirect("/auth/signin")
	}

	docId := c.Params("docId")

	if docId == "" {
		c.Set("HX-Redirect", "/dashboard")
		return c.Status(fiber.StatusNotFound).Redirect("/dashboard")
	}

	doc, err := database.GetDocumentFromId(docId)
	if err != nil {
		c.Set("HX-Redirect", "/dashboard/root")
		return c.Status(fiber.StatusNotFound).Redirect("/dashboard")
	}

	if doc.Type != "document" {
		c.Set("HX-Redirect", "/dashboard/root")
		return c.Status(fiber.StatusNotFound).Redirect("/dashboard")
	}

	page := documentView(user, *doc)
	handler := adaptor.HTTPHandler(templ.Handler(page))

	return handler(c)
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

func renameModal(c *fiber.Ctx) error {
	docId := c.Params("docId")

	doc, err := database.GetDocumentFromId(docId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	html := `
		<div class="modal modal-open" id="renameModal">
			<div class="modal-box">
				<h3 class="font-bold text-lg">Rename Document</h3>
				<form hx-post="/document/rename/` + docId + `"
					hx-target="#title-` + docId + `"
					hx-swap="innerHTML"
					hx-on="htmx:afterRequest:window.closeModal()"
					class="flex flex-col"
				>
					<input type="text" name="newName" class="input input-bordered mt-4" placeholder="New document name" value="` + doc.Title + `" required>
					<div class="modal-action">
						<a href="#" class="btn" onclick="closeModal()">Cancel</a>
						<button type="submit" class="btn btn-primary">Rename</button>
					</div>
				</form>
			</div>
		</div>

		<script>
			function closeModal() {
				const modal = document.getElementById('renameModal');
				if (modal) {
					modal.remove();
				}
			}

			// Fermer la modale lorsqu'on clique en dehors
			window.addEventListener('click', function(event) {
				const modal = document.getElementById('renameModal');
				if (event.target === modal) {
					closeModal();
				}
			});
		</script>
	`

	return c.SendString(html)
}

func renameDocument(c *fiber.Ctx) error {
	docId := c.Params("docId")
	userId := c.Locals("userID").(string)
	newName := c.FormValue("newName")

	if newName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("New name cannot be empty")
	}

	// get the document to rename
	doc, err := database.GetDocumentFromId(docId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// check if the user is the author of the document
	if doc.UserId != userId {
		return c.Status(fiber.StatusUnauthorized).SendString("You are not the author of this document")
	}

	err = database.RenameDocumentById(docId, newName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(newName)
}
