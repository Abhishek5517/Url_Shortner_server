package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB(dbURL string) {
	connectionString := dbURL

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		log.Fatal("unable to parse config", err)
	}

	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("error connecting to database", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	log.Println("connected to database!!!")
}
