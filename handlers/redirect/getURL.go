package redirectUrl

import (
	"context"
	"fmt"
	"time"
	"urlShortner/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetActualURL(key string) (string, error) {

	actual_url, hits, createdAt, err := getUrlFromDB(context.Background(), key)

	if err != nil {

		return "", fmt.Errorf("url invalid")
	}

	if !isURLExpired(createdAt) {

		return "", fmt.Errorf("url expired")
	}
	err = updateHitCounts(context.Background(), key, hits)

	if err != nil {

		return "", fmt.Errorf("url invalid")
	}
	if err != pgx.ErrNoRows && actual_url != "" {

		return actual_url, nil
	}

	return "", fmt.Errorf("url invalid")

}

func getUrlFromDB(ctx context.Context, key string) (string, pgtype.Int8, time.Time, error) {
	var actual_url string
	var hits pgtype.Int8
	var createdAt time.Time
	query := "SELECT actual_url , hits , created_at FROM short_urls WHERE short_code = $1"

	err := database.DB.QueryRow(ctx, query, key).Scan(&actual_url, &hits, &createdAt)

	return actual_url, hits, createdAt, err
}

func updateHitCounts(ctx context.Context, key string, hits pgtype.Int8) error {

	query := "UPDATE short_urls SET hits = $1 WHERE short_code = $2"

	_, err := database.DB.Exec(ctx, query, hits.Int64+int64(1), key)

	return err

}

func isURLExpired(createdAt time.Time) bool {

	expiredAt := createdAt.Add(24 * 7 * time.Hour)

	currentTime := time.Now()

	return currentTime.UnixMilli() < expiredAt.UnixMilli()

}
