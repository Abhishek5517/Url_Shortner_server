package createUrl

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"urlShortner/database"
	"urlShortner/models"

	"github.com/jackc/pgx/v5"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func GenerateHash(url string, user *models.LoginClaims) (string, error) {

	// combinedString := url + user.Email.String + "urlShortner@1.0.0"

	// hash := sha256.Sum256([]byte(combinedString))

	// key := hex.EncodeToString(hash[:][:7])
	// key = key[:7]
	key, err := RandomString(7)

	if err != nil {
		return "", fmt.Errorf("error generating key")
	}

	err = updateMapURL(key, url, user)

	if err != nil {
		return "", fmt.Errorf("db issue")
	}
	return key, nil

}

func updateMapURL(key string, actualURL string, user *models.LoginClaims) error {

	if !checkUrlPresent(context.Background(), key) {
		return nil
	}

	err := insertNewUrl(context.Background(), user.Email.String, key, actualURL)

	if err != nil {
		return err
	}

	return nil

}

func insertNewUrl(ctx context.Context, email string, short_code string, actual_url string) error {

	query := "INSERT INTO short_urls ( email , short_code , actual_url ) VALUES ($1 , $2, $3)"

	_, err := database.DB.Exec(ctx, query, email, short_code, actual_url)

	if err != nil {

		return fmt.Errorf("error inserting url data")
	}
	return nil
}

func checkUrlPresent(ctx context.Context, key string) bool {
	var actual_url string
	query := "SELECT actual_url FROM short_urls WHERE short_code = $1"

	err := database.DB.QueryRow(ctx, query, key).Scan(&actual_url)

	return err == pgx.ErrNoRows
}
