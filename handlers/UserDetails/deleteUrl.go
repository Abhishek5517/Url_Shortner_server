package user

import (
	"context"
	"fmt"
	"urlShortner/database"
	"urlShortner/models"

	"github.com/gofiber/fiber/v2"
)

func DeleteUrl() fiber.Handler {

	return func(c *fiber.Ctx) error {

		key := c.Params("key")

		claims := c.Locals("claims").(*models.LoginClaims)

		err := DeleteUrlwithKey(context.Background(), claims.Email.String, key)

		if err != nil {
			return fmt.Errorf("error deleting the record")
		}
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "successfully deleted the URL",
		})

	}

}

func DeleteUrlwithKey(ctx context.Context, email string, key string) error {

	query := "DELETE FROM short_urls WHERE email = $1 AND short_code = $2"

	_, err := database.DB.Exec(ctx, query, email, key)

	return err
}
