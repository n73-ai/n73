package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
		database.CreateLog("messages", projectID, err.Error())
		return c.SendStatus(500)
	}

	switch payload.Type {
	case "text":
		id := uuid.NewString()
		err := database.CreateMessage(id, projectID, "assistant", payload.Text, model, 0, false, 0.0)
		if err != nil {
			database.CreateLog("messages", projectID, err.Error())
			return c.SendStatus(500)
		}
		SendToUser(projectID, id)

	case "result":
		id := uuid.NewString()
		err := database.CreateMessage(id, projectID, "metadata", "", model, payload.Duration, payload.IsError, payload.TotalCostUsd)
		if err != nil {
			database.CreateLog("messages", projectID, err.Error())
			return c.SendStatus(500)
		}

		err = database.UpdateProjectSessionID(projectID, payload.SessionID)
		if err != nil {
			database.CreateLog("projects", projectID, err.Error())
			return c.SendStatus(500)
		}

		SendToUser(projectID, id)

		err = database.UpdateProjectStatus(projectID, "Deploying")
		if err != nil {
			database.CreateLog("projects", projectID, err.Error())
			return c.SendStatus(500)
		}

		SendToUser(projectID, "Deploying")

		project, err := database.GetProjectByID(projectID)
		if err != nil {
			database.CreateLog("projects", projectID, err.Error())
			return c.SendStatus(500)
		}

		projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

		err = utils.TryBuildProject(project.ID)
		if err != nil {
			wsFormatError := fmt.Sprintf("build-error: %s", err.Error())
			SendToUser(projectID, wsFormatError)
			database.CreateLog("projects", projectID, err.Error())
			return c.SendStatus(500)
		}

		err = utils.CopyProjectToExisitingProject(project.ID)
		if err != nil {
			SendToUser(projectID, "Error")
			database.CreateLog("projects", projectID, err.Error())
			return c.SendStatus(500)
		}

		err = utils.GhPush(projectPath)
		if err != nil {
			database.CreateLog("projects", project.ID, err.Error())
			err := database.UpdateProjectStatus(project.ID, "Gh-Error")
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
			}
			SendToUser(projectID, "Error")
			return nil
		}

		err = utils.CfPush(project.Slug, projectPath)
		if err != nil {
			database.CreateLog("projects", project.ID, err.Error())

			err := database.UpdateProjectStatus(project.ID, "Error")
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
			}

			SendToUser(projectID, "Error")
			return nil
		}

		err = database.UpdateProjectStatus(projectID, "Deployed")
		if err != nil {
			database.CreateLog("projects", project.ID, err.Error())
			return nil
		}

		SendToUser(projectID, "Deployed")

	default:
		log.Println("Unknown type:", payload.Type)
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}

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
