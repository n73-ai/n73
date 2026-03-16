package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"encoding/base64"
	"fmt"
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
		Type          string  `json:"type"`
		IsError       bool    `json:"is_error"`
		Text          string  `json:"text"`
		Duration      int     `json:"duration"`
		TotalCostUsd  float64 `json:"total_cost_usd"`
		SessionID     string  `json:"session_id"`
		File          string  `json:"file"`
		BuildError    bool    `json:"build_error"`
		BuildErrorMsg string  `json:"build_error_msg"`
		Image         string  `json:"image"`
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

		// start image
		var base64Image string

		// Handle both "data:application/zip;base64,XXXX" and raw base64
		if strings.Contains(payload.File, "data:") {
			parts := strings.SplitN(payload.File, ",", 2)
			if len(parts) < 2 {
				fmt.Println("image err.")
			}
			base64Image = parts[1]
		} else {
			base64Image = payload.Image
		}

		// Decode base64
		imageData, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			fmt.Println("image err: ", err.Error())
		}

		// Create temp zip file
		tempImageDataPath := filepath.Join("/tmp", "screenshot.png")
		err = os.WriteFile(tempImageDataPath, imageData, 0644)
		if err != nil {
			fmt.Println("image err: ", err.Error())
		}
		defer os.Remove(tempImageDataPath) // Clean up after
		// end image

		// Update status to idle immediately so the UI stops showing "Thinking"
		err = database.UpdateProjectStatus(projectID, "idle")
		if err != nil {
			database.CreateLog("projects", projectID, err.Error())
			return nil
		}
		SendToUser(projectID, "idle")

		// here starts go func()
		go func() {

			if payload.BuildError {
				err := database.UpdateProjectErrorMsg(projectID, payload.BuildErrorMsg)
				if err != nil {
					fmt.Println("error UpdateProjectErrorMsg(): ", err.Error())
				}
				error2fix := fmt.Sprintf("build-error: %s", payload.BuildErrorMsg)
				SendToUser(projectID, error2fix)

			} else {

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
					fmt.Println("write file err: ", err.Error())
					database.CreateLog("projects", projectID, "Failed to write zip file: "+err.Error())
					return
				}
				defer os.Remove(tempZipPath)

				// Unzip into projects directory
				repositoriesPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)
				err = utils.Unzip(tempZipPath, repositoriesPath)
				if err != nil {
					fmt.Println("unzip err: ", err.Error())
					database.CreateLog("projects", projectID, "Unzip failed: "+err.Error())
					return
				}

				// Purge Bunny CDN cache so changes are visible immediately
				project, err := database.GetProjectByID(projectID)
				if err == nil && project.PullZoneID != "" {
					if err := utils.PurgePullZoneCache(project.PullZoneID); err != nil {
						fmt.Println("cache purge err: ", err.Error())
					}
				}

				// Signal the frontend to reload the iframe now that CDN is updated
				SendToUser(projectID, "deployed")
			}

		}()

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
