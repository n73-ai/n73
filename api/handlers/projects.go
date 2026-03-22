package handlers

import (
	"ai-zustack/database"
	"ai-zustack/fly"
	"ai-zustack/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddCustomDomain(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func PublishProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")

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

	projectsDir := filepath.Join(os.Getenv("ROOT_PATH"), "projects")
	err = utils.GhClone(project.GhRepo, projectsDir, project.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

	err = utils.NpmRunBuild(projectPath + "/ui-only")
	if err != nil {
		fmt.Println("npm err: ", err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "npm run build err",
		})
	}

	if project.BunnyStatus == "storage_zone" {
		// name, region string
		mainRegion := "SE"
		storageZoneID, storageZonePassword, err := utils.CreateStorageZone(project.ID, mainRegion)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// here update project.bunny_region
		err = database.UpdateProjectStorageZone(project.ID, storageZoneID, storageZonePassword, "upload", "se")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// zonePassword, storageZoneName, distPath, region string
		distDir := filepath.Join(projectPath, "ui-only", "dist")
		err = utils.UploadDirectory(storageZonePassword, project.ID, distDir, project.StorageZoneRegion)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "upload")
			if err2 != nil {
				// log the error
				fmt.Println(err2.Error())
			}
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := utils.CreatePullZone(storageZoneID, project.ID, eu, na, asia, sa, af)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "pull_zone")
			if err2 != nil {
				// log the error
				fmt.Println(err2.Error())
			}
			// update the project.bunny_status == 'create_pull_zone'
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectPullZoneID(project.ID, pullZoneID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectDomain(domain, project.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = utils.DeleteProjectDirectory(projectPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(200)

	} else if project.BunnyStatus == "upload" {

		// zonePassword, storageZoneName, distPath, region string
		distDir := filepath.Join(projectPath, "ui-only", "dist")
		err = utils.UploadDirectory(project.StorageZonePassword, project.ID, distDir, project.StorageZoneRegion)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "upload")
			if err2 != nil {
				// log the error
				fmt.Println(err2.Error())
			}
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := utils.CreatePullZone(project.StorageZoneID, project.ID, eu, na, asia, sa, af)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "pull_zone")
			if err2 != nil {
				// log the error
				fmt.Println(err2.Error())
			}
			// update the project.bunny_status == 'create_pull_zone'
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectPullZoneID(project.ID, pullZoneID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectDomain(domain, project.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = utils.DeleteProjectDirectory(projectPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(200)

	} else if project.BunnyStatus == "pull_zone" {

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := utils.CreatePullZone(project.StorageZoneID, project.ID, eu, na, asia, sa, af)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "pull_zone")
			if err2 != nil {
				// log the error
				fmt.Println(err2.Error())
			}
			// update the project.bunny_status == 'create_pull_zone'
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectPullZoneID(project.ID, pullZoneID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = database.UpdateProjectDomain(domain, project.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = utils.DeleteProjectDirectory(projectPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(200)

	} else if project.BunnyStatus == "success" {
		err := utils.DeleteAllFilesInStorageZone(project.StorageZonePassword, project.StorageZoneID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = utils.UploadDirectory(project.StorageZonePassword, project.ID, projectPath, project.StorageZoneRegion)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = utils.PurgePullZoneCache(project.PullZoneID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

	} else {
		fmt.Println("Bunny Status unknow")
		return c.Status(500).JSON(fiber.Map{
			"error": "Bunny Status unknow",
		})
	}

	return c.SendStatus(200)
}

func TransferProject(c *fiber.Ctx) error {
	projectID := c.Params("projectID")
	email := c.Params("email")
	user := c.Locals("user").(*database.User)

	newOwner, err := database.GetUserByEmail(email)
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

	err = database.UpdateProjectOwner(projectID, newOwner.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func AdminGetProjects(c *fiber.Ctx) error {
	projects, err := database.GetProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func GetAllDeployedProjects(c *fiber.Ctx) error {
	projects, err := database.GetDeployedProjects()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func UpdateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
	payload := struct {
		Name string `json:"name"`
	}{}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The name is required.",
		})
	}

	if len(payload.Name) > 255 {
		return c.Status(400).JSON(fiber.Map{
			"error": "The name should not have more than 55 characters.",
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

	err = database.UpdateProjectName(projectID, payload.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func DeleteProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")

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

	go func() {
		err := fly.DeleteApp(projectID)
		if err != nil {
			fmt.Println("delete app: ", err.Error())
		}

		if project.PullZoneID != "" {
			err = utils.DeletePullZone(project.PullZoneID)
			if err != nil {
				fmt.Println("delete pull zone: ", err.Error())
			}
		}

		if project.StorageZoneID != "" {
			err = utils.DeleteStorageZone(project.StorageZoneID)
			if err != nil {
				fmt.Println("delete storage zone: ", err.Error())
			}
		}
	}()

	//projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)

	/*
	if _, err := os.Stat(projectPath); err != nil {
		if os.IsNotExist(err) {
      err = utils.DeleteProjectDirectory(projectPath)
      if err != nil {
        return c.Status(500).JSON(fiber.Map{
          "error": err.Error(),
        })
      }
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error checking project directory: " + err.Error(),
		})
	}
	*/

	err = database.DeleteProject(projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func CreateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)

	payload := struct {
		Prompt string `json:"prompt"`
		Name   string `json:"name"`
		Model  string `json:"model"`
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

	if len(payload.Model) > 55 {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	if payload.Model == "" {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	validModels := map[string]bool{
		"claude-sonnet-4-5-20250929": true,
		"claude-sonnet-4-20250514":   true,
		"claude-haiku-4-5-20251001":  true,
	}

	if !validModels[payload.Model] {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	hasIncident, err := database.HasActiveMajorFlyioIncident()
	if err == nil && hasIncident {
		return c.Status(503).JSON(fiber.Map{
			"error": "New projects can't be deployed right now due to a Fly.io incident. Check status.fly.io for updates.",
		})
	}

	projectID := uuid.NewString()
	payload.Name = "project-" + projectID
	messageID := uuid.NewString()
	webhookURL := fmt.Sprintf("%s/webhook/messages/%s/%s", os.Getenv("DOMAIN"), projectID, payload.Model)
	projectPath := filepath.Join("/app", "ui-only")

	slug := strings.ToLower(strings.ReplaceAll(payload.Name, " ", "-"))

	err = database.CreateProject(projectID, user.ID, payload.Name, slug)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateMessage(messageID, projectID, "user", payload.Prompt, payload.Model, 0, false, 0.0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	go func() {

		flyHostname, err := fly.CreateApp(projectID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = database.UpdateFlyHostname(projectID, flyHostname)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = fly.GenerateFlyToml(projectID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fileName := fmt.Sprintf("%s.toml", projectID)
		flyTomlPath := filepath.Join(os.Getenv("ROOT_PATH"), "fly_configs", fileName)
		err = fly.CreateMachine(flyTomlPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mainRegion := "SE"
		storageZoneID, storageZonePassword, err := utils.CreateStorageZone(projectID, mainRegion)
		if err != nil {
			fmt.Println("create storage zone: ", err.Error())
			return
		}

		err = database.UpdateProjectStorageZone(projectID, storageZoneID, storageZonePassword, "upload", "se")
		if err != nil {
			fmt.Println("update storage zone: ", err.Error())
			return
		}

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := utils.CreatePullZone(storageZoneID, projectID, eu, na, asia, sa, af)
		if err != nil {
			fmt.Println("create pull zone: ", err.Error())
			return
		}

		err = database.UpdateProjectPullZoneID(projectID, pullZoneID)
		if err != nil {
			fmt.Println("update pull zone id: ", err.Error())
			return
		}

		err = database.UpdateProjectDomain(domain, projectID)
		if err != nil {
			fmt.Println("update domain: ", err.Error())
			return
		}

		err = database.UpdateBunnyStatus(projectID, "success")
		if err != nil {
			fmt.Println("update bunny status: ", err.Error())
			return
		}

		endpoint := fmt.Sprintf("https://%s.fly.dev/claude/new", projectID)

		time.Sleep(5 * time.Second)
		err = utils.CreateClaudeProject(payload.Prompt, payload.Model, webhookURL, projectPath, endpoint, projectID, storageZonePassword)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}()

	return c.Status(200).JSON(fiber.Map{
		"project_id": projectID,
	})
}

func ResumeProject(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
	payload := struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
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

	if len(payload.Model) > 55 {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	if payload.Model == "" {
		payload.Model = "claude-sonnet-4-5-20250929"
	}

	validModels := map[string]bool{
		"claude-sonnet-4-5-20250929": true,
		"claude-sonnet-4-20250514":   true,
		"claude-haiku-4-5-20251001":  true,
	}

	if !validModels[payload.Model] {
		payload.Model = "claude-sonnet-4-5-20250929"
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

	if project.Status == "new_error" || project.Status == "new_internal_error" {
    err := database.UpdateProjectStatus(projectID, "new_pending")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if project.Status == "error" ||
		project.Status == "internal_error" ||
		project.Status == "idle" {
    err := database.UpdateProjectStatus(projectID, "pending")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	err = database.UpdateProjectErrorMsg(project.ID, "") 
	if err != nil {
		fmt.Println("error UpdateProjectErrorMsg(): ", err.Error())
	}

	go func() {

		messageID := uuid.NewString()
		endpoint := fmt.Sprintf("https://%s.fly.dev/claude/resume", projectID)
		webhookURL := fmt.Sprintf("%s/webhook/messages/%s/%s", os.Getenv("DOMAIN"), projectID, payload.Model)
		sessionID := project.SessionID

    err := utils.ResumeClaudeProject(
			payload.Prompt,
			payload.Model,
			webhookURL,
			"/app/ui-only",
			sessionID,
			endpoint,
			projectID,
			project.StorageZonePassword,
		)
		if err != nil {
      fmt.Println("resume claude project: ", err.Error())
			database.UpdateProjectStatus(projectID, "idle")
			SendToUser(projectID, "error")
			return
		}

    fmt.Println("projectID: ", projectID)
    fmt.Println("project.ID: ", project.ID)
    // err create message:  pq: insert or update on table "messages" violates foreign key constraint "messages_project_id_fkey"
		err = database.CreateMessage(messageID, project.ID, "user", payload.Prompt, payload.Model, 0, false, 0.0)
		if err != nil {
      fmt.Println("err create message: ", err.Error())
			SendToUser(projectID, "error")
			database.UpdateProjectStatus(projectID, "idle")
			return
		}
	}()

	return c.SendStatus(200)
}

func GetUserProjects(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projects, err := database.GetProjectsByUserID(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(projects)
}

func GetProjectByID(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	projectID := c.Params("projectID")
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
	return c.JSON(project)
}
