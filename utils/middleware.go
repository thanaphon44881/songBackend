package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"os"

	"github.com/joho/godotenv"
)

func AuthRequired(c *fiber.Ctx) error {
	godotenv.Load()
	Name := os.Getenv("NCOOKIE")
	cookie := c.Cookies(Name)
	jwtKey := os.Getenv("JWTKEY")
	if cookie == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("username", claims.Username)

	return c.Next()
}
