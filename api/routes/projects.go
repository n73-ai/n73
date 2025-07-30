package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func ProjectRoutes(app *fiber.App) {
	app.Post("/projects", handlers.CreateProject)
}
