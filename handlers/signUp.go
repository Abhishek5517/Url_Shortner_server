package auth

// name
// email
// password
// role

import (
	"context"

	"urlShortner/database"
	"urlShortner/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (pgtype.Text, error) {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return pgtype.Text{}, err
	}

	return pgtype.Text{String: string(hashedPass)}, nil
}

// add user data in DB / dummy signUpData.json file
func SignUp() fiber.Handler {
	var user models.User
	return func(c *fiber.Ctx) error {

		err := c.BodyParser(&user)

		if handleError(err, c) != nil {
			return handleError(err, c)
		}

		user.Password, err = HashPassword(user.Password.String)

		if handleError(err, c) != nil {
			return handleError(err, c)
		}

		var userData []models.User

		_, err = GetUserByEmail(context.Background(), database.DB, user.Email.String)

		if err != pgx.ErrNoRows {
			return c.Status(409).JSON(fiber.Map{
				"success":   false,
				"message":   "User already exists with this email",
				"errorCode": "USER_ALREADY_EXISTS",
			})
		}

		user.Id = pgtype.Int4{Int32: int32(len(userData) + 1)}

		err = insertUser(context.Background(), user)

		if handleError(err, c) != nil {
			return handleError(err, c)
		}

		token, err := IssueToken(&user)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success":   false,
				"message":   "Something went wrong. Please try again later.",
				"errorCode": "INTERNAL_SERVER_ERROR",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"message": "User registered successfully",
			"data": fiber.Map{
				"token": fiber.Map{
					"accessToken": token,
					"expiresIn":   3600,
				},
			},
		},
		)

	}

}

func handleError(err error, c *fiber.Ctx) error {

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success":   false,
			"message":   "Something went wrong. Please try again later.",
			"errorCode": "INTERNAL_SERVER_ERROR",
		})
	}
	return nil
}

func insertUser(ctx context.Context, user models.User) error {

	query := "INSERT INTO users ( name , email , password_hash ) VALUES ( $1 , $2 , $3 )"
	_, err := database.DB.Exec(ctx, query, user.Name.String, user.Email.String, user.Password.String)
	return err
}

func GetUserByEmail(ctx context.Context, DB *pgxpool.Pool, email string) (models.User, error) {
	var user models.User
	query := "SELECT name , email , password_hash FROM users WHERE email = $1; "

	err := DB.QueryRow(ctx, query, email).Scan(&user.Name,
		&user.Email, &user.Password)

	return user, err
}
