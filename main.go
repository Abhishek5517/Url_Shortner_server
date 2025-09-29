package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"urlShortner/database"
	auth "urlShortner/handlers"
	user "urlShortner/handlers/UserDetails"
	"urlShortner/handlers/createUrl"
	redirectUrl "urlShortner/handlers/redirect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// var opt, _ = redis.ParseURL(os.Getenv("REDIS_URL"))
// var rdb = redis.NewClient(opt)

func main() {

	_ = godotenv.Load()

	redis := database.ConnectRedis()

	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("SERVER_PORT")

	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	database.ConnectDB(dbURL)
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Referrer-Policy", "no-referrer")
		return c.Next()
	})

	app.Post("/signUp", auth.SignUp())

	app.Post("/login", auth.Login())
	app.Get("/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")

		pubsub := redis.Subscribe(context.Background(), "hits_channel")
		ch := pubsub.Channel()

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			for msg := range ch {
				fmt.Fprintf(w, "data: %s\n\n", msg.Payload)
				w.Flush()
			}
		})

		return nil
	})

	app.Get("/:key", redirectUrl.RedirectUrl())

	app.Use(auth.JWTauth())

	app.Post("/user", user.GetUserDetail())

	app.Post("/create", createUrl.CreateUrl())

	app.Delete("/delete/:key", user.DeleteUrl())

	log.Fatal(app.Listen(":" + port))

}
