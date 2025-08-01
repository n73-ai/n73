package handlers

import (
	"ai-zustack/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMessagesByProjectID(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
  messages, err := database.GetMessagesByProjectID(projectID)
  if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
  }
  return c.JSON(messages)
}

func WebhookMessage(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	metadataID := c.Params("metadataID")

	payload := struct {
		Type         string  `json:"type"`
		IsError      bool    `json:"is_error"`
		Text         string  `json:"text"`
		Duration     int     `json:"duration"`
		TotalCostUsd float64 `json:"total_cost_usd"`
		SessionID    string  `json:"session_id"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	switch payload.Type {
	case "text":
	  id := uuid.NewString()
		err := database.CreateMessage(id, projectID, metadataID, "assistant", payload.Text)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
    SendToUser("user_id", id)
	case "result":
		err := database.UpdateMetadata(payload.SessionID, payload.Duration, payload.TotalCostUsd, metadataID, payload.IsError)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = database.UpdateProjectSessionID(projectID, payload.SessionID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
    SendToUser("user_id", "done")
	default:
		log.Println("Unknown type:", payload.Type)
	}

	return c.SendStatus(200)
}

