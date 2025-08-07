package api

import (
	"ai-zustack/api/handlers"
	"ai-zustack/api/routes"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RunServer() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))

	app.Use("/feed/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	go handlers.RunHub()

	routes.ProjectsRoutes(app)
	routes.MessagesRoutes(app)
	routes.WebsocketRoutes(app)
	routes.UsersRoutes(app)

	return app
}
