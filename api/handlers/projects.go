package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

	// hard coded for now
	userID := "42069"
	projectID := uuid.NewString()
	projectName := "Full-Stack-App"
	model := "claude-3-5-haiku-20241022"
	messageID := uuid.NewString()

	webhookURL := fmt.Sprintf("http://0.0.0.0:%s/webhook/messages/%s/%s", os.Getenv("PORT"), projectID, model)

	err = utils.CreateClaudeProject(payload.Prompt, model, webhookURL)
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
