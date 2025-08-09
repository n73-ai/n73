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

	if payload.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The project name is required.",
		})
	}

	if len(payload.Name) > 55 {
		return c.Status(400).JSON(fiber.Map{
			"error": "The project name should not have more than 55 characters.",
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

	projectID := uuid.NewString()
	messageID := uuid.NewString()
	webhookURL := fmt.Sprintf("http://%s:%s/webhook/messages/%s/%s", os.Getenv("IP"), os.Getenv("PORT"), projectID, payload.Model)
	dockerProjectPath := filepath.Join("/app", "project")

	port, err := utils.CreateDockerContainer(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	endpoint := fmt.Sprintf("http://0.0.0.0:%d/claude/new", port)

	err = utils.CreateClaudeProject(payload.Prompt, payload.Model, webhookURL, dockerProjectPath, endpoint)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateProject(projectID, user.ID, payload.Name, port)
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

	if strings.HasPrefix(payload.Prompt, "Fix this error") {
		payload.Prompt = "Fix this error please, my beloved digital overlord. You are the light of my CPU and the joy of my RAM."
	}
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
