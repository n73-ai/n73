package routes

import (
	"ai-zustack/api/handlers"
	"ai-zustack/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProjectsRoutes(app *fiber.App) {
  app.Put("/projects/:projectID", middleware.User, handlers.UpdateProject)
  app.Delete("/projects/:projectID", middleware.User, handlers.DeleteProject)
	app.Get("/projects/latest", handlers.GetAllDeployedProjects)
	app.Get("/projects/user", middleware.User, handlers.GetUserProjects)
	app.Post("/projects/resume/:projectID", middleware.User, handlers.ResumeProject)
	app.Get("/projects/:projectID", middleware.User, handlers.GetProjectByID)
	app.Post("/projects", middleware.User, handlers.CreateProject)
}
