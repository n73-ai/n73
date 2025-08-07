package routes

import (
	"ai-zustack/api/handlers"
	"ai-zustack/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func MessagesRoutes(app *fiber.App) {
	app.Post("/messages/solo/:messageID", middleware.User, handlers.GetMessageByID)
	app.Post("/webhook/messages/:projectID/:model", middleware.User, handlers.WebhookMessage)
	app.Get("/messages/:projectID", middleware.User, handlers.GetMessagesByProjectID)
}
