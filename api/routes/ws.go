package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebsocketRoutes(app *fiber.App) {
	app.Get("/feed/chat", websocket.New(handlers.FeedChat))
}
