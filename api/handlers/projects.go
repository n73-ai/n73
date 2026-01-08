package handlers

import (
	"ai-zustack/database"
	"ai-zustack/fly"
	"ai-zustack/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PublishProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")

	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if project.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You don't have access to this resource",
		})
	}

	slug := strings.ToLower(strings.ReplaceAll(project.Name, " ", "-"))

	cloudflareProjectName := fmt.Sprintf("project-%s", projectID)
	exists, err := utils.PageExists(cloudflareProjectName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !exists {
		err = utils.CfCreate(slug)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	domain, err := utils.GetProjectDomainFallback(slug)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.UpdateProjectDomain(projectID, domain)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	projectsDir := filepath.Join(os.Getenv("ROOT_PATH"), "projects")
	err = utils.GhClone(project.GhRepo, projectsDir, project.ID)
	if err != nil {
		fmt.Println("gh error:", err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "GitHub clone error",
		})
	}

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

	err = utils.NpmRunBuild(projectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "npm run build err",
		})
	}

	err = utils.CfPush(project.Slug, projectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Deploy error to Cloudflare",
		})
	}

	// delete the directory to free disk
	err = utils.DeleteProjectDirectory(projectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func TransferProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	email := c.Params("email")
	user := c.Locals("user").(*database.User)

	newOwner, err := database.GetUserByEmail(email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if project.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You don't have access to this resource",
		})
	}

	err = database.UpdateProjectOwner(projectID, newOwner.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func AdminGetProjects(c *fiber.Ctx) error {
	projects, err := database.GetProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func GetAllDeployedProjects(c *fiber.Ctx) error {
	projects, err := database.GetDeployedProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func UpdateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
	payload := struct {
		Name string `json:"name"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The name is required.",
		})
	}

	if len(payload.Name) > 255 {
		return c.Status(400).JSON(fiber.Map{
			"error": "The name should not have more than 55 characters.",
		})
	}

	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if project.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You don't have access to this resource",
		})
	}

	err = database.UpdateProjectName(projectID, payload.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func DeleteProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")

	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if project.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You don't have access to this resource",
		})
	}

	go func() {
		err := fly.DeleteApp(projectID)
		if err != nil {
			fmt.Println("delete app: ", err.Error())
		}
		err = utils.DeleteGhRepo(projectID)
		if err != nil {
			fmt.Println("delete gh repo: ", err.Error())
		}
		err = utils.DeleteCfPage(projectID)
		if err != nil {
			fmt.Println("cf page err: ", err.Error())
		}
	}()

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)

	if _, err := os.Stat(projectPath); err != nil {
		if os.IsNotExist(err) {
			err = database.DeleteProject(projectID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.SendStatus(200)
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error checking project directory: " + err.Error(),
		})
	}

	err = os.RemoveAll(projectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error removing project directory: " + err.Error(),
		})
	}

	err = database.DeleteProject(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func CreateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)

	payload := struct {
		Prompt string `json:"prompt"`
		Name   string `json:"name"`
		Model  string `json:"model"`
	}{}

	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Prompt == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The prompt is required.",
		})
	}

	if len(payload.Model) > 55 {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	if payload.Model == "" {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	validModels := map[string]bool{
		"claude-sonnet-4-5-20250929": true,
		"claude-sonnet-4-20250514":   true,
		"claude-haiku-4-5-20251001":  true,
	}

	if !validModels[payload.Model] {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	projectID := uuid.NewString()
	payload.Name = "project-" + projectID
	messageID := uuid.NewString()
	webhookURL := fmt.Sprintf("%s/webhook/messages/%s/%s", os.Getenv("DOMAIN"), projectID, payload.Model)
	projectPath := filepath.Join("/app", "ui-only")

	slug := strings.ToLower(strings.ReplaceAll(payload.Name, " ", "-"))

	err = database.CreateProject(projectID, user.ID, payload.Name, slug)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, payload.Model, 0, false, 0.0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	go func() {

		flyHostname, err := fly.CreateApp(projectID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = database.UpdateFlyHostname(projectID, flyHostname)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = fly.GenerateFlyToml(projectID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fileName := fmt.Sprintf("%s.toml", projectID)
		flyTomlPath := filepath.Join(os.Getenv("ROOT_PATH"), "fly_configs", fileName)
		err = fly.CreateMachine(flyTomlPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		endpoint := fmt.Sprintf("http://%s.internal:5000/claude/new", projectID)

		time.Sleep(10 * time.Second)
		err = utils.CreateClaudeProject(payload.Prompt, payload.Model, webhookURL, projectPath, endpoint)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)
		err = utils.GhCreate(slug, projectPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ghRepo := fmt.Sprintf("https://github.com/n73-projects/%s", slug)
		err = database.UpdateGhRepo(projectID, ghRepo)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}()

	return c.Status(200).JSON(fiber.Map{
		"project_id": projectID,
	})
}

func ResumeProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
	payload := struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Prompt == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Prompt can't be empty.",
		})
	}

	if len(payload.Model) > 55 {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	if payload.Model == "" {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	validModels := map[string]bool{
		"claude-sonnet-4-5-20250929": true,
		"claude-sonnet-4-20250514":   true,
		"claude-haiku-4-5-20251001":  true,
	}

	if !validModels[payload.Model] {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if project.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You don't have access to this resource",
		})
	}

	if project.Status == "new_error" || project.Status == "new_internal_error" {
		err = database.UpdateProjectStatus(projectID, "new_pending")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if project.Status == "error" ||
		project.Status == "internal_error" ||
		project.Status == "idle" {
		err = database.UpdateProjectStatus(projectID, "pending")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	go func() {

		messageID := uuid.NewString()
		endpoint := fmt.Sprintf("http://%s.internal:5000/claude/new", projectID)
		webhookURL := fmt.Sprintf("%s/webhook/messages/%s/%s", os.Getenv("DOMAIN"), projectID, payload.Model)
		sessionID := project.SessionID

		err = utils.ResumeClaudeProject(
			payload.Prompt,
			payload.Model,
			webhookURL,
			"/app/ui-only",
			sessionID,
			endpoint)
		if err != nil {
			database.UpdateProjectStatus(projectID, "idle")
			SendToUser(projectID, "error")
			return
		}

		err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, payload.Model, 0, false, 0.0)
		if err != nil {
			SendToUser(projectID, "error")
			database.UpdateProjectStatus(projectID, "idle")
			return
		}
	}()

	return c.SendStatus(200)
}

func GetUserProjects(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projects, err := database.GetProjectsByUserID(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func GetProjectByID(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if project.UserID != user.ID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You don't have access to this resource",
		})
	}
	return c.JSON(project)
}
