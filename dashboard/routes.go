package dashboard

import (
	"time"

	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/ZiplEix/Google-Docs-Wish/middleware"
	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/dashboard", middleware.Protected())
	dashboardGroup.Get("/", dashboardPage)
	dashboardGroup.Get("/search", dashboardSearch)
}

func dashboardSearch(c *fiber.Ctx) error {
	query := c.Query("q")

	type result struct {
		Title  string
		Link   string
		Date   time.Time
		Author string
	}

	results := []result{
		{
			Title:  query,
			Link:   "https://www.google.com",
			Date:   time.Now(),
			Author: "Client",
		},
		{
			Title:  "Test",
			Link:   "https://www.google.com",
			Date:   time.Now(),
			Author: "Author 1",
		},
		{
			Title:  "Test 2",
			Link:   "https://www.google.com",
			Date:   time.Now(),
			Author: "Author 2",
		},
		{
			Title:  "Test 3",
			Link:   "https://www.google.com",
			Date:   time.Now(),
			Author: "Author 3",
		},
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
            <a href="` + res.Link + `" class="block p-4 border-b border-base-300 hover:bg-base-200 transition-colors duration-200">
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
                        ` + res.Date.Format("January 2, 2006") + `
                    </p>
                </div>
            </a>
        `
	}

	return c.SendString(html)
}

func generateDocumentListHtml(user database.User) string {
	documents, err := database.GetDocumentFromUserId(user.ID)
	if err != nil {
		return ""
	}

	var html string

	html += "<div class='w-5/6'>"

	for _, doc := range documents {
		html += `
			<div class="flex items-center w-full h-auto px-4 py-3 rounded-3xl bg-base-200 mb-4 shadow-sm pl-6">
				<!-- Document Info -->
				<a href="/document/` + doc.ID + `" class="flex-1">
					<p class="text-lg font-semibold">` + doc.Title + `</p>
					<p class="text-sm text-gray-500">Last modified on ` + time.Now().Format("2 January 2006") + `</p>
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
						<li><a href="#" onclick="deleteDocument('` + doc.ID + `')">Delete Document</a></li>
						<li><a href="#" onclick="renameDocument('` + doc.ID + `')">Rename Document</a></li>
					</ul>
				</div>
			</div>
		`
	}

	html += "</div>"

	return html
}
