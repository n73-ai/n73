package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
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

		project, err := database.GetProjectByID(projectID)
		if err != nil {
			database.CreateLog("projects", projectID, err.Error())
			return c.SendStatus(500)
		}

		projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

		err = utils.TryBuildProject(project.ID)
		if err != nil {

			errP := database.UpdateProjectStatus(project.ID, "new_error")
			if errP != nil {
				database.CreateLog("projects", project.ID, err.Error())
			}

			errP = database.UpdateProjectErrorMsg(project.ID, err.Error())
			if errP != nil {
				database.CreateLog("projects", project.ID, err.Error())
			}

			isFirstBuild := project.Status == "new_pending"
			if isFirstBuild {
				SendToUser(projectID, "new_error")
			} else {
				SendToUser(projectID, "error")
			}

			return c.SendStatus(500)
		}

		isFirstBuild := project.Status == "new_pending"
		var pStatusErr string
		if isFirstBuild {
			pStatusErr = "new_internal_error"
		} else {
			pStatusErr = "internal_error"
		}

		err = utils.CopyProjectToExisitingProject(project.ID)
		if err != nil {
			database.CreateLog("Copy Project error", project.ID, err.Error())

			err := database.UpdateProjectStatus(project.ID, pStatusErr)
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
			}

			SendToUser(projectID, "error")
			return c.SendStatus(500)
		}

		err = utils.CfPush(project.Slug, projectPath)
		if err != nil {
			database.CreateLog("Cloudflare push error", project.ID, err.Error())

			err := database.UpdateProjectStatus(project.ID, pStatusErr)
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
			}

			SendToUser(projectID, "error")
			return nil
		}

		err = utils.GhPush(projectPath)
		if err != nil {
			database.CreateLog("GitHub push error", project.ID, err.Error())
		}

		err = database.UpdateProjectStatus(projectID, "idle")
		if err != nil {
			database.CreateLog("projects", project.ID, err.Error())
			return nil
		}

		SendToUser(projectID, "idle")

	default:
		log.Println("Unknown type:", payload.Type)
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}

func GetMessageByID(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	messageID := c.Params("messageID")
	message, err := database.GetMessageByID(messageID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	project, err := database.GetProjectByID(message.ProjectID)
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
	return c.JSON(message)
}

func GetMessagesByProjectID(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
	messages, err := database.GetMessagesByProjectID(projectID)
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
	return c.JSON(messages)
}
