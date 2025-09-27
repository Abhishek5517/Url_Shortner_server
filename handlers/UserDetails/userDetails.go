package user

import (
	"context"
	"urlShortner/database"
	"urlShortner/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetUserDetail() fiber.Handler {

	return func(c *fiber.Ctx) error {

		claims := c.Locals("claims").(*models.LoginClaims)
		email := claims.Email.String

		// if email == "" {
		// 	return fmt.Errorf("no email provided for user details")
		// }
		rows, err := getUserFromDB(context.Background(), email)

		if err != nil {
			return err
		}
		var mappedUrls []models.UrlMap

		for rows.Next() {
			var urlmap models.UrlMap

			rows.Scan(&urlmap.Short_code, &urlmap.Actual_url, &urlmap.CreatedAt, &urlmap.Hits)

			mappedUrls = append(mappedUrls, urlmap)
		}

		if len(mappedUrls) == 0 {
			return c.Status(200).JSON(fiber.Map{
				"success": true,
				"message": "no urls found",
				"data":    []string{},
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"data":    mappedUrls,
		})

	}

}

func getUserFromDB(ctx context.Context, email string) (pgx.Rows, error) {

	query := "SELECT short_code , actual_url, created_at , hits FROM short_urls WHERE email = $1"

	rows, err := database.DB.Query(ctx, query, email)

	return rows, err

}
