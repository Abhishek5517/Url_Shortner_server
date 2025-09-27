package auth

import (
	"context"
	"fmt"
	"os"
	"time"
	"urlShortner/database"
	"urlShortner/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwt secret
var secretKey = os.Getenv("SECRET_KEY")

var JwtSecret = []byte(secretKey)

func authenticateUser(user *models.User) (string, error) {

	if user.Email.String == "" {
		return "", fmt.Errorf("invalid email")
	}

	signedUpUser, err := GetUserByEmail(context.Background(), database.DB, user.Email.String)

	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("user not signed up please sign up ")
	}

	if bcrypt.CompareHashAndPassword([]byte(signedUpUser.Password.String), []byte(user.Password.String)) != nil {
		return "", fmt.Errorf("invalid password")
	} else {
		return signedUpUser.Password.String, nil
	}

}

func IssueToken(user *models.User) (string, error) {

	claims := models.LoginClaims{
		Email:    user.Email,
		Password: user.Password,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "2bit.com",
			Audience:  []string{"2bit-clients"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-30 * time.Second)),
			ID:        "jti-" + fmt.Sprint(time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)

}

func Login() fiber.Handler {

	return func(c *fiber.Ctx) error {

		user := new(models.User)

		err := c.BodyParser(&user)

		if err != nil {
			return handleInvalidation(c)
		}

		pass, err := authenticateUser(user)
		if err != nil {
			return handleInvalidation(c)
		}
		if pass != "" {

			token, err := IssueToken(user)

			if err != nil {
				return handleInvalidation(c)
			}

			return c.Status(200).JSON(fiber.Map{
				"success": true,
				"message": "Login successful",
				"data": fiber.Map{
					"token": fiber.Map{
						"accessToken": token,
						"expiresIn":   3600,
					},
				},
			})
		}

		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "invalid credentials",
		})
	}

}

func handleInvalidation(c *fiber.Ctx) error {

	return c.Status(400).JSON(fiber.Map{
		"success":   false,
		"message":   "Email and password are required",
		"errorCode": "VALIDATION_ERROR",
		"details": fiber.Map{
			"email":    "Email is required",
			"password": "Password is required",
		},
	})
}
