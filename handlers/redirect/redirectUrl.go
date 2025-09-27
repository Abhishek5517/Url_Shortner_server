package redirectUrl

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var RedisUrl = os.Getenv("REDIS_URL")
var rdb = redis.NewClient(&redis.Options{
	Addr: RedisUrl,
})

var ctx = context.Background()

func RedirectUrl() fiber.Handler {

	return func(c *fiber.Ctx) error {

		key := c.Params("key")

		actualURL, err := GetActualURL(key)

		if err != nil {
			return c.SendFile("./public/index.html")
		}

		count, err := rdb.Incr(ctx, "url:"+key).Result()
		if err != nil {
			return c.Status(500).SendString("Redis error")
		}

		rdb.Publish(ctx, "hits_channel", fmt.Sprintf("%s:%d", key, count))
		return c.Redirect(actualURL)

	}
}
