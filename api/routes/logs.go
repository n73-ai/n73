package routes

import (
	"ai-zustack/api/handlers"
	"ai-zustack/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func LogsRoutes(app *fiber.App) {
	app.Get("/logs", middleware.Admin, handlers.GetLogs)
}
