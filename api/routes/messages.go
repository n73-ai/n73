package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func MessagesRoutes(app *fiber.App) {
	app.Post("/messages/solo/:messageID", handlers.GetMessageByID)
	app.Post("/webhook/messages/:projectID/:model", handlers.WebhookMessage)
	app.Get("/messages/:projectID", handlers.GetMessagesByProjectID)
}
