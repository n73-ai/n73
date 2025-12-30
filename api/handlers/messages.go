package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"

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
		File         string  `json:"file"`
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

		// here starts go func()
		go func() {

			// === HANDLE BASE64 ZIP FILE ===
			if payload.File == "" {
			  database.CreateLog("projects", projectID, "missing base64 zip file")
        return
			}

			var base64Data string

			// Handle both "data:application/zip;base64,XXXX" and raw base64
			if strings.Contains(payload.File, "data:") {
				parts := strings.SplitN(payload.File, ",", 2)
				if len(parts) < 2 {
			    database.CreateLog("projects", projectID, "Invalid data URL format")
          return
				}
				base64Data = parts[1]
			} else {
				base64Data = payload.File
			}

			// Decode base64
			zipData, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				database.CreateLog("projects", projectID, "Base64 decode error: "+err.Error())
				return 
			}

			// Create temp zip file
			tempZipPath := filepath.Join("/tmp", projectID+".zip")
			err = os.WriteFile(tempZipPath, zipData, 0644)
			if err != nil {
				database.CreateLog("projects", projectID, "Failed to write zip file: "+err.Error())
				return 
			}
			defer os.Remove(tempZipPath) // Clean up after

			// Unzip into projects directory
			repositoriesPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)
			err = utils.Unzip(tempZipPath, repositoriesPath)
			if err != nil {
				database.CreateLog("projects", projectID, "Unzip failed: "+err.Error())
				return 
			}

			projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)
			// frontendProjectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID, "project")

			/*
			    err = utils.TryBuildProject(frontendProjectPath)
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
			*/

			/*
				isFirstBuild := project.Status == "new_pending"
				var pStatusErr string
				if isFirstBuild {
					pStatusErr = "new_internal_error"
				} else {
					pStatusErr = "internal_error"
				}

				err = utils.CfPush(project.Slug, frontendProjectPath)
				if err != nil {
					database.CreateLog("Cloudflare push error", project.ID, err.Error())

					err := database.UpdateProjectStatus(project.ID, pStatusErr)
					if err != nil {
						database.CreateLog("projects", project.ID, err.Error())
					}

					SendToUser(projectID, "error")
					return nil
				}
			*/

			err = utils.GhPush(projectPath)
			if err != nil {
				database.CreateLog("GitHub push error", project.ID, err.Error())
			}
			// here ends go func()
		}()

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
