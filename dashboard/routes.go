package dashboard

import (
	"time"

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
