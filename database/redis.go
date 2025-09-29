package database

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis() redis.Client {

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("ParseURL failed: %v", err)
	}

	opt.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	return *redis.NewClient(opt)

}
