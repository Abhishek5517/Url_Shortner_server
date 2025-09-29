package redirectUrl

import (
	"context"
	"fmt"
	"urlShortner/database"

	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

func RedirectUrl() fiber.Handler {

	redis := database.ConnectRedis()

	return func(c *fiber.Ctx) error {

		key := c.Params("key")

		actualURL, err := GetActualURL(key)

		if err != nil {
			return c.SendFile("./public/index.html")
		}

		count, err := redis.Incr(ctx, "url:"+key).Result()
		if err != nil {
			return c.Status(500).SendString("Redis error")
		}

		redis.Publish(ctx, "hits_channel", fmt.Sprintf("%s:%d", key, count))
		return c.Redirect(actualURL)

	}
}
