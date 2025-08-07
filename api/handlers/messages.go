package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		SendToUser("hej@agustfricke.com", id)

		err = database.UpdateProjectStatus(projectID, "Deploying")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		SendToUser("hej@agustfricke.com", "deploy-start")

		project, err := database.GetProjectByID(projectID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)
		err = utils.NpmRunBuild(projectPath)
		if err != nil {
			wsFormatError := fmt.Sprintf("build-error: %s", err.Error())
			SendToUser("hej@agustfricke.com", wsFormatError)
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if project.CFProjectReady {
			err = utils.Push(project.Name, projectPath)
			if err != nil {
				updateProjectError := database.UpdateProjectStatus(projectID, "Failed deployment")
				if updateProjectError != nil {
					return c.Status(500).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateProjectStatus(projectID, "Deployed")
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			SendToUser("hej@agustfricke.com", "deploy-done")
		} else {
			err = utils.FistDeployment(project.Name, projectPath)
			if err != nil {
				updateProjectError := database.UpdateProjectStatus(projectID, "Failed deployment")
				if updateProjectError != nil {
					return c.Status(500).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateProjectCFProjectReady(projectID, true)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateProjectStatus(projectID, "Deployed")
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			slug := strings.ToLower(strings.ReplaceAll(project.Name, " ", "-"))
			domain, err := utils.GetProjectDomainFallback(slug)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateProjectDomain(projectID, domain)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			SendToUser("hej@agustfricke.com", "deploy-done")
		}
	default:
		log.Println("Unknown type:", payload.Type)
	}
	return c.SendStatus(200)
}
