package handlers

import (
	"ai-zustack/database"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMessageByID(c *fiber.Ctx) error {
  messageID := c.Params("messageID")
  message, err := database.GetMessageByID(messageID)
  if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
  }
  return c.JSON(message)
}

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
	model := c.Params("model")

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
		err := database.CreateMessage(id, projectID, "assistant", payload.Text, model, 0, false, 0.0)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		SendToUser("hej@agustfricke.com", id)
	case "result":
		id := uuid.NewString()
		err := database.CreateMessage(id, projectID, "metadata", "", model, payload.Duration, payload.IsError, payload.TotalCostUsd)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectSessionID(projectID, payload.SessionID)
		if err != nil {
      fmt.Println("err: ", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var status string
		if payload.IsError {
			status = "AI-Error"
		} else {
			status = "Ready"
		}

		err = database.UpdateProjectStatus(projectID, status)
		if err != nil {
      fmt.Println("err: ", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		SendToUser("hej@agustfricke.com", id)

	default:
		log.Println("Unknown type:", payload.Type)
	}

	return c.SendStatus(200)
}
