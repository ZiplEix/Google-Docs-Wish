package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("auth_token")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.Set("HX-Redirect", "/auth/signin")
			return c.Status(fiber.StatusUnauthorized).Redirect("/auth/signin")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					c.Set("HX-Redirect", "/auth/signin")
					return c.Status(fiber.StatusUnauthorized).Redirect("/auth/signin")
				}
			}

			userID := claims["sub"].(string)
			c.Locals("userID", userID)

			return c.Next()
		}

		c.Set("HX-Redirect", "/dashboard")
		return c.Status(fiber.StatusUnauthorized).Redirect("/dashboard")
	}
}
