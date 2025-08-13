package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func KeepClaudeAlive(prompt, endpoint string) error {
	payload := map[string]string{
		"prompt": prompt,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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

func ResumeClaudeProject(prompt, model, webhookURL, path, sessionID, endpoint string) error {
	payload := map[string]string{
		"work_dir":    path,
		"prompt":      prompt,
		"model":       model,
		"webhook_url": webhookURL,
		"session_id":  sessionID,
		"jwt":         os.Getenv("ADMIN_JWT"),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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

func CreateClaudeProject(prompt, model, webhookURL, path, endpoint string) error {
	payload := map[string]string{
		"work_dir":    path,
		"prompt":      prompt,
		"model":       model,
		"webhook_url": webhookURL,
		"jwt":         os.Getenv("ADMIN_JWT"),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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
