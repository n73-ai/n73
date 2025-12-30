package routes

import (
	"ai-zustack/api/handlers"
	"ai-zustack/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProjectsRoutes(app *fiber.App) {
	app.Get("/admin/projects", middleware.Admin, handlers.AdminGetProjects)

	app.Post("/projects/transfer/:projectID/:email", middleware.User, handlers.TransferProject)

	app.Put("/projects/:projectID", middleware.User, handlers.UpdateProject)
	app.Delete("/projects/:projectID", middleware.User, handlers.DeleteProject)

	app.Get("/projects/latest", handlers.GetAllDeployedProjects)
	app.Get("/projects/user", middleware.User, handlers.GetUserProjects)
	app.Post("/projects/resume/:projectID", middleware.User, handlers.ResumeProject)
	app.Get("/projects/:projectID", middleware.User, handlers.GetProjectByID)
	app.Post("/projects", middleware.User, handlers.CreateProject)
}
