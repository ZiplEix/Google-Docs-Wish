package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/ZiplEix/Google-Docs-Wish/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitFirestore()
	database.InitFirebaseAuth()

	// auth.InitOauth()
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Static("/", "./public")

	app.Use(logger.New())

	handler.Setup(app)

	fmt.Println("Server is running on http://localhost:" + os.Getenv("PORT"))
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
