package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func MessagesRoutes(app *fiber.App) {
  app.Post("/webhook/messages/:projectID/:metadataID", handlers.WebhookMessage)
  app.Get("/messages/:projectID", handlers.GetMessagesByProjectID)
}
