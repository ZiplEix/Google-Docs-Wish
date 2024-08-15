package auth

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/ZiplEix/Google-Docs-Wish/database"
	"github.com/ZiplEix/Google-Docs-Wish/users"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func isValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	return len(password) >= 8
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func setTokenCookie(c *fiber.Ctx, user users.User) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
		"sub": user.ID,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return nil
}

func signin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		fmt.Println("Email and password are required")
		return c.Status(fiber.StatusUnauthorized).SendString("Email and password are required")
	}

	if !isValidEmail(email) {
		fmt.Println("Invalid email")
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid email")
	}

	if !isValidPassword(password) {
		fmt.Println("Password must be at least 8 characters long")
		return c.Status(fiber.StatusUnauthorized).SendString("Password must be at least 8 characters long")
	}

	docsIter := database.FirestoreClient.Collection("users").Where("email", "==", email).Documents(c.Context())
	docs, err := docsIter.GetAll()
	if err != nil {
		fmt.Println("Error getting user data:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	if len(docs) == 0 {
		fmt.Println("User does not exist")
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	if len(docs) > 1 {
		fmt.Println("Multiple users with the same email")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	user := users.New(docs[0].Data(), docs[0].Ref.ID)

	err = verifyPassword(user.Password, password)
	if err != nil {
		fmt.Println("Invalid password")
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid password")
	}

	err = setTokenCookie(c, *user)
	if err != nil {
		fmt.Println("Error setting token cookie:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

func signup(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm-password")

	if email == "" || password == "" || confirmPassword == "" {
		fmt.Println("Email and password are required")
		return c.Status(fiber.StatusUnauthorized).SendString("Email and password are required")
	}

	if !isValidEmail(email) {
		fmt.Println("Invalid email")
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid email")
	}

	if !isValidPassword(password) {
		fmt.Println("Password must be at least 8 characters long")
		return c.Status(fiber.StatusUnauthorized).SendString("Password must be at least 8 characters long")
	}

	if password != confirmPassword {
		fmt.Println("Passwords do not match")
		return c.Status(fiber.StatusUnauthorized).SendString("Passwords do not match")
	}

	docsIter := database.FirestoreClient.Collection("users").Where("email", "==", email).Documents(c.Context())
	docs, err := docsIter.GetAll()
	if err != nil {
		fmt.Println("Error getting user data:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	if len(docs) > 0 {
		fmt.Println("User already exist")
		return c.Status(fiber.StatusUnauthorized).SendString("User already exist, email is already in use")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	user := users.New(map[string]interface{}{
		"email":    email,
		"password": hashedPassword,
	})

	_, err = user.CreateInDb()
	if err != nil {
		fmt.Println("Error creating user:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	err = setTokenCookie(c, *user)
	if err != nil {
		fmt.Println("Error setting token cookie:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error, please try again later")
	}

	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}
