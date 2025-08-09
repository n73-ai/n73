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
			fmt.Println(err.Error())
		}

		err = database.UpdateProjectSessionID(projectID, payload.SessionID)
		if err != nil {
			fmt.Println(err.Error())
		}

		SendToUser(projectID, id)

		err = database.UpdateProjectStatus(projectID, "Deploying")
		if err != nil {
			fmt.Println(err.Error())
		}

		SendToUser(projectID, "Deploying")

		project, err := database.GetProjectByID(projectID)
		if err != nil {
			fmt.Println(err.Error())
		}

		projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

		err = utils.TryBuildProject(project.ID)
		if err != nil {
			fmt.Println(err.Error())
			wsFormatError := fmt.Sprintf("build-error: %s", err.Error())
			SendToUser(projectID, wsFormatError)
		}

		// fist deployment, on every stage update p.stage & if err p.err_stage
		if project.Stage == "0" {
			// copy in new project
			err = utils.CopyProjectToMainMachine(project.ID)
			if err != nil {
				fmt.Println(err.Error())
				wsFormatError := fmt.Sprintf("build-error: %s", err.Error())
				SendToUser(projectID, wsFormatError)
			}

			// create github remote repository and push first code
			err = utils.GhCreate(project.Slug, projectPath)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}

			/*
				err = database.UpdateProjectStage(project.ID, "1")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
			*/

			// push to github repository
			err = utils.GhPush(projectPath)
			if err != nil {
				fmt.Println("a", err.Error())
				database.CreateLog("projects", project.ID, err.Error())
				err = database.UpdateProjectErrorStage(project.ID, "1")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
				return nil
			}
			/*
				err = database.UpdateProjectStage(project.ID, "2")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
			*/

			// create cf page
			err = utils.CfCreate(project.Slug)
			if err != nil {
				fmt.Println("b", err.Error())
				database.CreateLog("projects", project.ID, err.Error())
				err = database.UpdateProjectErrorStage(project.ID, "2")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
				return nil
			}
			/*
				err = database.UpdateProjectStage(project.ID, "3")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
			*/

			// push to cloudflare
			err = utils.CfPush(project.Slug, projectPath)
			if err != nil {
				fmt.Println("c", err.Error())
				fmt.Println(err.Error())
				database.CreateLog("projects", project.ID, err.Error())
				err = database.UpdateProjectErrorStage(project.ID, "3")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
				return nil
			}
			err = database.UpdateProjectStage(project.ID, "4")
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
				return nil
			}
			// finish deployment

			err = database.UpdateProjectStatus(projectID, "Deployed")
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
				return nil
			}

			domain, err := utils.GetProjectDomainFallback(project.Slug)
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
				return nil
			}
			err = database.UpdateProjectDomain(projectID, domain)
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
				return nil
			}

			SendToUser(projectID, "Deployed")
			return nil
		}

		if project.Stage == "4" {
			err = utils.CopyProjectToExisitingProject(project.ID)
			if err != nil {
				fmt.Println(err.Error())
				wsFormatError := fmt.Sprintf("build-error: %s", err.Error())
				SendToUser(projectID, wsFormatError)
			}
			// push to github repository
			err = utils.GhPush(projectPath)
			if err != nil {
				fmt.Println(err.Error())
				database.CreateLog("projects", project.ID, err.Error())
				err = database.UpdateProjectErrorStage(project.ID, "1")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
				return nil
			}

			/*
				err = database.UpdateProjectStage(project.ID, "2")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
			*/

			// push to cloudflare
			err = utils.CfPush(project.Slug, projectPath)
			if err != nil {
				fmt.Println(err.Error())
				database.CreateLog("projects", project.ID, err.Error())
				err = database.UpdateProjectErrorStage(project.ID, "3")
				if err != nil {
					database.CreateLog("projects", project.ID, err.Error())
					return nil
				}
				return nil
			}
			err = database.UpdateProjectStatus(projectID, "Deployed")
			if err != nil {
				database.CreateLog("projects", project.ID, err.Error())
				return nil
			}
			SendToUser(projectID, "Deployed")
			return nil
		}

	default:
		log.Println("Unknown type:", payload.Type)
		return c.SendStatus(400)
	}
	return c.SendStatus(200)
}
