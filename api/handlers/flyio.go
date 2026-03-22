package handlers

import (
	"ai-zustack/database"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type flyioIncidentPayload struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	Impact     string  `json:"impact"`
	ResolvedAt *string `json:"resolved_at"`
}

type flyioComponentUpdatePayload struct {
	Status    string `json:"status"`
	Component struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"component"`
}

type flyioWebhookPayload struct {
	Incident        *flyioIncidentPayload        `json:"incident"`
	ComponentUpdate *flyioComponentUpdatePayload `json:"component_update"`
}

func broadcastFlyioStatus() {
	incidents, _ := database.GetActiveFlyioIncidents()
	if incidents == nil {
		incidents = []database.FlyioIncident{}
	}
	msg, _ := json.Marshal(fiber.Map{
		"type":      "flyio_status",
		"incidents": incidents,
	})
	broadcast <- string(msg)
}

func FlyioWebhook(c *fiber.Ctx) error {
	var payload flyioWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if payload.Incident != nil {
		inc := payload.Incident
		resolved := inc.ResolvedAt != nil || inc.Status == "resolved" || inc.Status == "postmortem"
		if err := database.UpsertFlyioIncident(inc.ID, inc.Name, inc.Status, inc.Impact, resolved); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go broadcastFlyioStatus()
	}

	if payload.ComponentUpdate != nil {
		cu := payload.ComponentUpdate
		resolved := cu.Status == "operational"
		name := cu.Component.Name + " - " + cu.Status
		if err := database.UpsertFlyioIncident("component-"+cu.Component.ID, name, cu.Status, "", resolved); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go broadcastFlyioStatus()
	}

	return c.SendStatus(200)
}

func GetFlyioStatus(c *fiber.Ctx) error {
	incidents, err := database.GetActiveFlyioIncidents()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if incidents == nil {
		incidents = []database.FlyioIncident{}
	}
	return c.JSON(fiber.Map{"incidents": incidents})
}
