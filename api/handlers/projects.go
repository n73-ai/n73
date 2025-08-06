package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func DeployProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

	if project.CFProjectReady {
		err = utils.PushGH(projectPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = utils.PushCF(project.Name, projectPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(200)
	}

	err = utils.CreatePushGH(project.Name, projectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = utils.CreateCFPage(project.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = utils.PushCF(project.Name, projectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func CreateProject(c *fiber.Ctx) error {
	payload := struct {
		Prompt string `json:"prompt"`
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

	// hard coded
	userID := "42069"
	projectID := uuid.NewString()
	projectName := "Full-Stack-App"
	model := "claude-3-5-haiku-20241022"
	messageID := uuid.NewString()

	webhookURL := fmt.Sprintf("http://0.0.0.0:%s/webhook/messages/%s/%s", os.Getenv("PORT"), projectID, model)

	path := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	err = utils.CreateClaudeProject(payload.Prompt, model, webhookURL, path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateProject(projectID, userID, projectName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, model, 0, false, 0.0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"project_id": projectID,
	})
}
