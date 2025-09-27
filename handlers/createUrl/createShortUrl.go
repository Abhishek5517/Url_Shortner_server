package createUrl

import (
	"fmt"
	"urlShortner/models"

	"github.com/gofiber/fiber/v2"
)

type Url struct {
	URL string
}

func CreateUrl() fiber.Handler {

	return func(c *fiber.Ctx) error {

		claims := c.Locals("claims").(*models.LoginClaims)

		URL := new(Url)

		err := c.BodyParser(&URL)

		if err != nil {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"success": false,
				"message": "issue creating shorurl",
			})
		}

		key, err := GenerateHash(URL.URL, claims)

		if err != nil {
			return fmt.Errorf("db issues")
		}

		shortLink := "http://localhost:8888/" + key

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success":  true,
			"message":  "short url created successfully",
			"shortUrl": shortLink,
		})

	}
}
