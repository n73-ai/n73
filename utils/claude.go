package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func ResumeClaudeProject(prompt, model, webhookURL, path, sessionID string) error {
	payload := map[string]string{
		"work_dir":    path,
		"prompt":      prompt,
		"model":       model,
		"webhook_url": webhookURL,
		"session_id":  sessionID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := "http://localhost:5000/claude/resume"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: got %d (%s), expected 200 OK.", resp.StatusCode, resp.Status)
	}

	return nil
}

func CreateClaudeProject(prompt, model, webhookURL, path string) error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	payload := map[string]string{
		"work_dir":    path,
		"prompt":      prompt,
		"model":       model,
		"webhook_url": webhookURL,
		"api_key":     apiKey,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := "http://localhost:5000/claude/new"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: got %d (%s), expected 200 OK.", resp.StatusCode, resp.Status)
	}

	return nil
}
