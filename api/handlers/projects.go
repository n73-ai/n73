package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
		Name   string `json:"name"`
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
	projectID := c.Params("projectID")
	_, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = utils.RmDockerContainer(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)

	if _, err := os.Stat(projectPath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Project directory not found",
			})
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
		payload.Model = "claude-3-7-sonnet-20250219"
	}

	if payload.Model == "" {
		payload.Model = "claude-3-7-sonnet-20250219"
	}

	validModels := map[string]bool{
		"claude-sonnet-4-20250514":   true,
		"claude-3-7-sonnet-20250219": true,
		"claude-3-5-sonnet-20241022": true,
		"claude-3-5-sonnet-20240620": true,
		"claude-3-5-haiku-20241022":  true,
		"claude-3-haiku-20240307":    true,
	}

	if !validModels[payload.Model] {
		payload.Model = "claude-3-7-sonnet-20250219"
	}

	projectID := uuid.NewString()
	payload.Name = "project-" + projectID
	messageID := uuid.NewString()
	webhookURL := fmt.Sprintf("http://%s:%s/webhook/messages/%s/%s", os.Getenv("IP"), os.Getenv("PORT"), projectID, payload.Model)
	dockerProjectPath := filepath.Join("/app", "project")

	port, err := utils.GetFreePort()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	slug := strings.ToLower(strings.ReplaceAll(payload.Name, " ", "-"))
	err = database.CreateProject(projectID, user.ID, payload.Name, slug, port)
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

		err = utils.CreateDockerContainer(projectID, port)
		if err != nil {
      fmt.Println(err.Error())
			return
		}

		endpoint := fmt.Sprintf("http://0.0.0.0:%d/claude/new", port)

		err = utils.CreateClaudeProject(payload.Prompt, payload.Model, webhookURL, dockerProjectPath, endpoint)
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

		err = utils.CfCreate(slug)
		if err != nil {
      fmt.Println(err.Error())
			return
		}

		domain, err := utils.GetProjectDomainFallback(slug)
		if err != nil {
      fmt.Println(err.Error())
			return
		}

		err = database.UpdateProjectDomain(projectID, domain)
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

	if payload.Model == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The model is required.",
		})
	}

	validModels := map[string]bool{
		"claude-sonnet-4-20250514":   true,
		"claude-3-7-sonnet-20250219": true,
		"claude-3-5-sonnet-20241022": true,
		"claude-3-5-sonnet-20240620": true,
		"claude-3-5-haiku-20241022":  true,
		"claude-3-haiku-20240307":    true,
	}

	if !validModels[payload.Model] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid model. Please select a valid Claude model.",
		})
	}

	err = database.UpdateProjectStatus(projectID, "Building")
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

	messageID := uuid.NewString()
	webhookURL := fmt.Sprintf("http://%s:%s/webhook/messages/%s/%s", os.Getenv("IP"), os.Getenv("PORT"), projectID, payload.Model)
	dockerProjectPath := filepath.Join("/app", "project")
	endpoint := fmt.Sprintf("http://0.0.0.0:%d/claude/resume", project.Port)

	err = utils.ResumeClaudeProject(payload.Prompt, payload.Model, webhookURL, dockerProjectPath, project.SessionID, endpoint)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	/*
		if strings.HasPrefix(payload.Prompt, "Fix this build error") {
			payload.Prompt = "Fix this error please, my beloved digital overlord. You are the light of my CPU and the joy of my RAM."
		}
	*/
	err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, payload.Model, 0, false, 0.0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func GetUserProjects(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projects, err := database.GetProjectsByUserID(user.ID)
	if err != nil {
		fmt.Println("oo")
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func GetProjectByID(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(project)
}
