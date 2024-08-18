package dashboard

import (
	"fmt"

	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/ZiplEix/Google-Docs-Wish/middleware"
	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/dashboard", middleware.Protected())
	dashboardGroup.Get("/", dashboardPage)
	dashboardGroup.Get("/search", dashboardSearch)
	dashboardGroup.Get("/:rootId", dashboardPage)
}

func dashboardSearch(c *fiber.Ctx) error {
	query := c.Query("q")

	if query == "" {
		return c.SendString("")
	}

	userId := c.Locals("userID").(string)

	results, err := database.SearchDocument(query, userId)
	if err != nil {
		fmt.Printf("error searching documents: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var html string

	if len(results) == 0 {
		html = `
			<div class="text-center text-gray-500">
				No results found
			</div>
		`

		return c.SendString(html)
	}

	for _, res := range results {
		html += `
            <a href="/document/` + res.ID + `" class="block p-4 border-b border-base-300 hover:bg-base-200 transition-colors duration-200">
                <div class="flex justify-between items-start">
                    <div class="flex-1">
                        <div class="font-semibold text-lg text-blue-600">
                            ` + res.Title + `
                        </div>
                        <p class="text-sm text-gray-500">
                            ` + res.Author + `
                        </p>
                    </div>
                    <p class="text-sm text-gray-500">
                        ` + res.LastModified.Format("January 2, 2006") + `
                    </p>
                </div>
            </a>
        `
	}

	return c.SendString(html)
}

func generateDocumentLIstTile(doc *database.Document, image, url string) string {
	return `
			<div id="doc-` + doc.ID + `" class="flex items-center w-full h-auto px-4 py-3 rounded-3xl bg-base-200 mb-4 shadow-sm">
				<!-- Document Icon -->
				<img src="` + image + `" class="w-10 h-10 mr-2" alt="Document Icon">

				<!-- Document Info -->
				<a href="` + url + `" class="flex-1">
					<p class="text-lg font-semibold">` + doc.Title + `</p>
					<p class="text-sm text-gray-500">Last modified on ` + doc.LastModified.Format("2 January 2006") + `</p>
				</a>

				<!-- Dropdown Menu -->
				<div class="dropdown dropdown-left">
					<label tabindex="0" class="btn btn-circle btn-ghost btn-md">
						<svg fill="currentColor" width="800px" height="800px" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg" class="inline-block w-6 h-6">
							<path d="M9.5 13a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0zm0-5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0zm0-5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"/>
						</svg>
					</label>
					<ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-48">
						<li><a href="/document/` + doc.ID + `" target="_blank">Open in New Tab</a></li>
						<li><a href="#" hx-delete="/document/` + doc.ID + `" hx-target="#doc-` + doc.ID + `" hx-swap="outerHTML">Delete Document</a></li>
						<li><a href="#" onclick="renameDocument('` + doc.ID + `')">Rename Document</a></li>
					</ul>
				</div>
			</div>
		`
}

func generateBreadcrumbPath(rootId string) string {
	currentDirectory, err := database.GetDocumentFromId(rootId)
	if err != nil || currentDirectory == nil {
		return `<nav class='flex items-center text-sm text-gray-500'><a href='/dashboard' class='text-blue-500 hover:text-blue-700'>Home</a></nav>`
	}

	html := `<nav class="flex items-center text-base text-gray-500 mb-4">`

	path := `
		<span class="mx-2">/</span>` +
		`<a href="/dashboard/` + currentDirectory.ID + `" class="text-blue-500 hover:text-blue-700">` + currentDirectory.Title + `</a>
	`

	for currentDirectory.RootId != `root` {
		currentDirectory, err = database.GetDocumentFromId(currentDirectory.RootId)
		if err != nil || currentDirectory == nil {
			break
		}

		path = `<span class="mx-2">/</span>` +
			`<a href="/dashboard/` + currentDirectory.ID + `" class="text-blue-500 hover:text-blue-700">` + currentDirectory.Title + `</a>` +
			path
	}

	path = `<a href="/dashboard" class="text-blue-500 hover:text-blue-700">Home</a>` + path

	html += path + `</nav>`

	return html
}

func generateDocumentListHtml(user database.User, rootId string) string {
	documents, err := database.GetDocumentFromUserId(user.ID, rootId)
	if err != nil {
		return ""
	}

	var html string
	html += "<div class='w-5/6'>"

	if rootId != "root" {
		html += generateBreadcrumbPath(rootId)
	}

	for _, doc := range documents {
		image := "/ui/"

		switch doc.Type {
		case "directory":
			image += "directory_icon.png"
		case "document":
			image += "doc_icon.png"
		case "spreadsheet":
			image += "xls_icon.png"
		case "pdf":
			image += "pdf_icon.png"
		default:
			image += "doc_icon.png"
		}

		url := "/document/" + doc.ID
		if doc.Type == "directory" {
			url = "/dashboard/" + doc.ID
		}

		html += generateDocumentLIstTile(doc, image, url)
	}

	html += "</div>"
	return html
}
