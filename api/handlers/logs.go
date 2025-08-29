package handlers

import (
	"ai-zustack/database"

	"github.com/gofiber/fiber/v2"
)

func GetLogs(c *fiber.Ctx) error {
	logs, err := database.GetLogs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(logs)
}
