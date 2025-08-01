package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func CreateClaudeProject(prompt, model, webhookURL string) error {
	workDir := filepath.Join(os.Getenv("ROOT_PATH"), "ai-project")
	payload := map[string]string{
		"work_dir":    workDir,
		"prompt":      prompt,
		"model":       model,
		"webhook_url": webhookURL,
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
