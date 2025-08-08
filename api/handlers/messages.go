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
		SendToUser(projectID, id)

	case "result":
		id := uuid.NewString()
		err := database.CreateMessage(id, projectID, "metadata", "", model, payload.Duration, payload.IsError, payload.TotalCostUsd)
		if err != nil {
      fmt.Println("1", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectSessionID(projectID, payload.SessionID)
		if err != nil {
      fmt.Println("2", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		SendToUser(projectID, id)

		err = database.UpdateProjectStatus(projectID, "Deploying")
		if err != nil {
      fmt.Println("3", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		SendToUser(projectID, "deploy-start")

		project, err := database.GetProjectByID(projectID)
		if err != nil {
      fmt.Println("4", err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)
		err = utils.NpmRunBuild(projectPath)
		if err != nil {
      fmt.Println("5", err.Error())
			wsFormatError := fmt.Sprintf("build-error: %s", err.Error())
			SendToUser(projectID, wsFormatError)
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if project.CFProjectReady {
			err = utils.Push(project.Name, projectPath)
			if err != nil {
        fmt.Println("6", err.Error())
				updateProjectError := database.UpdateProjectStatus(projectID, "Failed deployment")
				if updateProjectError != nil {
          fmt.Println("7", updateProjectError.Error())
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
        fmt.Println("8", err.Error())
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			SendToUser(projectID, "deploy-done")
		} else {
			err = utils.FistDeployment(project.Name, projectPath)
			if err != nil {
        fmt.Println("9", err.Error())
				updateProjectError := database.UpdateProjectStatus(projectID, "Failed deployment")
				if updateProjectError != nil {
          fmt.Println("10", updateProjectError.Error())
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
        fmt.Println("11", err.Error())
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateProjectStatus(projectID, "Deployed")
			if err != nil {
        fmt.Println("12", err.Error())
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			slug := strings.ToLower(strings.ReplaceAll(project.Name, " ", "-"))
			domain, err := utils.GetProjectDomainFallback(slug)
      //project%20domains’
			if err != nil {
        fmt.Println("13", err.Error())
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateProjectDomain(projectID, domain)
			if err != nil {
        fmt.Println("14", err.Error())
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			SendToUser(projectID, "deploy-done")
		}
	default:
		log.Println("Unknown type:", payload.Type)
	}
	return c.SendStatus(200)
}
