package auth

import (
	"fmt"
	"strings"
	"urlShortner/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JWTauth() fiber.Handler {

	return func(c *fiber.Ctx) error {

		auth := c.Get("Authorization")

		if strings.HasPrefix(auth, "Bearer ") {
			auth = strings.TrimPrefix(auth, "Bearer ")
		} else {
			return c.Status(401).JSON(fiber.Map{
				"success":   false,
				"message":   "Token is expired or invalid. Please login again.",
				"errorCode": "TOKEN_EXPIRED",
			})
		}

		claims := &models.LoginClaims{}
		token, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {

			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("invalid signing methods")
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {

			return c.Status(401).JSON(fiber.Map{
				"success":   false,
				"message":   "Token is expired or invalid. Please login again.",
				"errorCode": "TOKEN_EXPIRED",
			})
		}
		c.Locals("claims", claims)
		return c.Next()
	}

}
