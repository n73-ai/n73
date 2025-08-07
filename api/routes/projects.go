package routes

import (
	"ai-zustack/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func ProjectsRoutes(app *fiber.App) {
	app.Get("/projects/user", handlers.GetUserProjects)
	app.Post("/projects/resume/:projectID", handlers.ResumeProject)
	app.Get("/projects/:projectID", handlers.GetProjectByID)
	app.Post("/projects/deploy/:projectID", handlers.DeployProject)
	app.Post("/projects", handlers.CreateProject)
}
