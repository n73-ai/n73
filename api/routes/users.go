package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UsersRoutes(app *fiber.App) {
	app.Post("/users/auth/link", handlers.AuthLink)
	app.Post("/users/auth/verify/:tokenString", handlers.VerifyAuthLink)
}
