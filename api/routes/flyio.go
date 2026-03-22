package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func FlyioRoutes(app *fiber.App) {
	app.Post("/webhook/flyio", handlers.FlyioWebhook)
	app.Get("/flyio/status", handlers.GetFlyioStatus)
}
