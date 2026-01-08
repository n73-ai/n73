package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func PurgePullZoneCache(pullZoneID int64) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("%s/pullzone/%d/purgeCache", apiBase, pullZoneID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error purging cache (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func AddRedirectEdgeRule(pullZoneID int64, bunnyHostname string, customDomain string) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	payload := map[string]any{
		"ActionType":          1,
		"ActionParameter1":    customDomain,
		"ActionParameter2":    "301",
		"ActionParameter3":    "",
		"Description":         "redirect bunny domain",
		"Enabled":             true,
		"ExtraActions":        []any{},
		"TriggerMatchingType": 0,
		"Triggers": []map[string]any{
			{
				"PatternMatches":      []string{bunnyHostname},
				"PatternMatchingType": 0,
				"Type":                0,
				"Parameter1":          "",
			},
		},
		"OrderIndex": 0,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	url := fmt.Sprintf("%s/pullzone/%d/edgerules/addOrUpdate", apiBase, pullZoneID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating edge rule (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func LoadFreeCertificate(customHostname string) error {
	const apiBase = "https://api.bunny.net/pullzone/loadFreeCertificate"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("%s?Hostname=%s", apiBase, customHostname)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error loading free SSL certificate (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func EnableForceSSL(pullZoneID int64, customHostname string) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	payload := map[string]any{
		"ForceSSL": true,
		"Hostname": customHostname,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	url := fmt.Sprintf("%s/pullzone/%d/setForceSSL", apiBase, pullZoneID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error enabling Force SSL (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func AddCustomHostname(pullZoneID int64, customHostname string) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	payload := map[string]string{
		"Hostname": customHostname,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling hostname: %w", err)
	}

	url := fmt.Sprintf("%s/pullzone/%d/addHostname", apiBase, pullZoneID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error adding hostname (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func CreatePullZone() error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	createPayload := map[string]any{
		"Name":              "kool-pullzone",
		"Type":              0,
		"OriginType":        2,
		"StorageZoneId":     1326407,
		"EnableGeoZoneEU":   true,
		"EnableGeoZoneUS":   false,
		"EnableGeoZoneASIA": false,
		"EnableGeoZoneSA":   false,
		"EnableGeoZoneAF":   false,
	}

	jsonData, err := json.Marshal(createPayload)
	if err != nil {
		return fmt.Errorf("error marshaling create payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiBase+"/pullzone", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating pull zone (status %d): %s", resp.StatusCode, string(body))
	}

	var pullZoneResp struct {
		ID int64 `json:"Id"`
	}
	if err := json.Unmarshal(body, &pullZoneResp); err != nil {
		return fmt.Errorf("error parsing respuesta: %w", err)
	}
	// pullZoneID := pullZoneResp.ID

	return nil
}

func UploadDirectory() error {
	const localDir = "/home/agust/dist"
	const storageZoneName = "kool-name"
	const region = "se"
	const baseURL = "https://" + region + ".storage.bunnycdn.com/" + storageZoneName + "/"
	const zonePassword = "af5ae962-dc68-48c9-af576a5ee87d-78ca-45cd"

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	err := filepath.Walk(localDir, func(localPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localDir, localPath)
		if err != nil {
			return fmt.Errorf("error calculation relative route: %w", err)
		}
		objectPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		url := baseURL + objectPath

		file, err := os.Open(localPath)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", localPath, err)
		}
		defer file.Close()

		req, err := http.NewRequest("PUT", url, file)
		if err != nil {
			return fmt.Errorf("error creating request for %s: %w", objectPath, err)
		}

		req.Header.Set("AccessKey", zonePassword)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error uploading %s: %w", objectPath, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("error uploading %s: status %d, body: %s", objectPath, resp.StatusCode, string(body))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error uploading directory: %w", err)
	}

	return nil
}

func CreateStorageZone(name string) error {
	url := "https://api.bunny.net/storagezone"

	payload := map[string]any{
		"Name":            name,
		"Region":          "SE",
		"StorageZoneType": 0,
		"ZoneTier":        "Standard",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling create payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AccessKey", os.Getenv("BUNNYNET_ACCESS_KEY"))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
