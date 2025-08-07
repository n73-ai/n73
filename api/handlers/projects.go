package handlers

import (
	"ai-zustack/database"
	"ai-zustack/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func DeployProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	project, err := database.GetProjectByID(projectID)
	if err != nil {
    fmt.Println("0")
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

	if project.CFProjectReady {
    // push code to gh
		err = utils.PushGH(projectPath)
		if err != nil {
      fmt.Println("222")
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
    // push code to cf
		err = utils.PushCF(project.Name, projectPath)
		if err != nil {
      fmt.Println("333")
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(200)
	}

  // creates new github repo & push the remote repository
	err = utils.CreatePushGH(project.Name, projectPath)
	if err != nil {
    fmt.Println("453")
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  // creates new cf page
	err = utils.CreateCFPage(project.Name)
	if err != nil {
    fmt.Println("1")
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  err = database.UpdateProjectCFProjectReady(projectID, true)
  if err != nil {
    fmt.Println("2")
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
  }

  // push the code under ./dist to cf pages
	err = utils.PushCF(project.Name, projectPath)
	if err != nil {
      fmt.Println("34")
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func GetUserProjects(c *fiber.Ctx) error {
  userID := "42069"
  projects, err := database.GetProjectsByUserID(userID)
  if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
  }
  return c.JSON(projects)
}

func ResumeProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	payload := struct {
		Prompt string `json:"prompt"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Prompt == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Prompt can't be empty.",
		})
	}

	err = database.UpdateProjectStatus(projectID, "Building")
	if err != nil {
    fmt.Println("err: ", err.Error())
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

	// hard coded
	model := "claude-3-5-haiku-20241022"
	messageID := uuid.NewString()

	webhookURL := fmt.Sprintf("http://0.0.0.0:%s/webhook/messages/%s/%s", os.Getenv("PORT"), projectID, model)

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)
  if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return c.Status(500).JSON(fiber.Map{
			"error": "project directory does not exist.",
		})
  } else if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
  }

	err = utils.ResumeClaudeProject(payload.Prompt, model, webhookURL, projectPath, project.SessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, model, 0, false, 0.0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  return c.SendStatus(200)
}

func GetProjectByID(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	project, err := database.GetProjectByID(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(project)
}


func CreateProject(c *fiber.Ctx) error {
	payload := struct {
		Prompt string `json:"prompt"`
		Name string `json:"name"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Prompt == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The prompt is required.",
		})
	}

	if payload.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The project name is required.",
		})
	}

	if len(payload.Name) > 55 {
		return c.Status(400).JSON(fiber.Map{
			"error": "The project name should not have more than 55 characters.",
		})
	}

	// !hard coded
	userID := "42069"
  model := "claude-3-5-haiku-20241022"
  // end!

	projectID := uuid.NewString()
	messageID := uuid.NewString()

	webhookURL := fmt.Sprintf("http://0.0.0.0:%s/webhook/messages/%s/%s", os.Getenv("PORT"), projectID, model)

	baseProjectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", "base-project")
	newProjectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)
	cmd := exec.Command("cp", "-r", baseProjectPath, newProjectPath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = utils.CreateClaudeProject(payload.Prompt, model, webhookURL, newProjectPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateProject(projectID, userID, payload.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, model, 0, false, 0.0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"project_id": projectID,
	})
}
