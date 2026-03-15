package utils

import (
	"ai-zustack/database"
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

// update bunny status correctly, after each {}, current status, but is success
func Upload2Bunny(project database.Project) error {
	projectPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", project.ID)

	if project.BunnyStatus == "storage_zone" {
		// name, region string
		mainRegion := "SE"
		storageZoneID, storageZonePassword, err := CreateStorageZone(project.ID, mainRegion)
		if err != nil {
			return err
		}

		// here update project.bunny_region
		err = database.UpdateProjectStorageZone(project.ID, storageZoneID, storageZonePassword, "upload", "se")
		if err != nil {
			return err
		}

		// zonePassword, storageZoneName, distPath, region string
		distDir := filepath.Join(projectPath, "dist")
		err = UploadDirectory(storageZonePassword, project.ID, distDir, project.StorageZoneRegion)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "upload")
			if err2 != nil {
				return err2
			}
			return err
		}

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := CreatePullZone(storageZoneID, project.ID, eu, na, asia, sa, af)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "pull_zone")
			if err2 != nil {
				return err2
			}
			return err
		}

		err = database.UpdateProjectPullZoneID(project.ID, pullZoneID)
		if err != nil {
			return err
		}

		err = database.UpdateProjectDomain(domain, project.ID)
		if err != nil {
			return err
		}

		err = database.UpdateBunnyStatus(project.ID, "success")
		if err != nil {
			return err
		}

    /*
		err = DeleteProjectDirectory(projectPath)
		if err != nil {
			return err
		}
    */

		return nil

	} else if project.BunnyStatus == "upload" {

		// zonePassword, storageZoneName, distPath, region string
		distDir := filepath.Join(projectPath, "dist")
		err := UploadDirectory(project.StorageZonePassword, project.ID, distDir, project.StorageZoneRegion)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "upload")
			if err2 != nil {
				return err2
			}
			return err
		}

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := CreatePullZone(project.StorageZoneID, project.ID, eu, na, asia, sa, af)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "pull_zone")
			if err2 != nil {
				return err2
			}
			return err
		}

		err = database.UpdateProjectPullZoneID(project.ID, pullZoneID)
		if err != nil {
			return err
		}

		err = database.UpdateProjectDomain(domain, project.ID)
		if err != nil {
			return err
		}

    /*
		err = DeleteProjectDirectory(projectPath)
		if err != nil {
			return err
		}
    */

		err = database.UpdateBunnyStatus(project.ID, "success")
		if err != nil {
			return err
		}

		return nil

	} else if project.BunnyStatus == "pull_zone" {

		eu := true
		na := false
		asia := false
		sa := false
		af := false
		pullZoneID, domain, err := CreatePullZone(project.StorageZoneID, project.ID, eu, na, asia, sa, af)
		if err != nil {
			err2 := database.UpdateBunnyStatus(project.ID, "pull_zone")
			if err2 != nil {
				return err2
			}
			return err
		}

		err = database.UpdateProjectPullZoneID(project.ID, pullZoneID)
		if err != nil {
			return err
		}

		err = database.UpdateProjectDomain(domain, project.ID)
		if err != nil {
			return err
		}

    /*
		err = DeleteProjectDirectory(projectPath)
		if err != nil {
			return err
		}
    */

		err = database.UpdateBunnyStatus(project.ID, "success")
		if err != nil {
			return err
		}

		return nil

	} else if project.BunnyStatus == "success" {
    fmt.Printf("storage zone id: %s, storage zone password: %s", project.StorageZoneID, project.StorageZonePassword)
    err := DeleteAllFilesInStorageZone(project.ID, project.StorageZonePassword)
		if err != nil {
			return err
		}

		distDir := filepath.Join(projectPath, "dist")
		err = UploadDirectory(project.StorageZonePassword, project.ID, distDir, project.StorageZoneRegion)
		if err != nil {
			return err
		}

		err = PurgePullZoneCache(project.PullZoneID)
		if err != nil {
			return err
		}

    return nil

	} 

  return fmt.Errorf("Error: project.bunny_status unknow.")

}

func DeleteAllFilesInStorageZone(storageZoneName, password string) error {
	region := "se"
	baseURL := fmt.Sprintf("https://%s.storage.bunnycdn.com/%s", region, storageZoneName)
	accessKey := password

	if accessKey == "" {
		return fmt.Errorf("password is empty (use your Storage Zone password)")
	}

	client := &http.Client{Timeout: 60 * time.Second}

	listURL := baseURL + "/"
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return fmt.Errorf("error creating GET list request: %w", err)
	}
	req.Header.Set("AccessKey", accessKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error listing files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error listing (status %d): %s", resp.StatusCode, string(body))
	}

	// Minimal type just for decoding (no full struct)
	var objects []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&objects); err != nil {
		return fmt.Errorf("error parsing JSON list: %w", err)
	}

	if len(objects) == 0 {
		fmt.Println("The Storage Zone is already empty.")
		return nil
	}

	// Step 2: Delete each object (files and folders)
	for _, obj := range objects {
		objectName, _ := obj["ObjectName"].(string)
		isDirectory, _ := obj["IsDirectory"].(bool)

		relativePath := objectName
		if isDirectory && !strings.HasSuffix(relativePath, "/") {
			relativePath += "/"
		}

		deleteURL := baseURL + "/" + relativePath

		delReq, err := http.NewRequest("DELETE", deleteURL, nil)
		if err != nil {
			fmt.Printf("Error creating DELETE request for %s: %v\n", relativePath, err)
			continue
		}
		delReq.Header.Set("AccessKey", accessKey)

		delResp, err := client.Do(delReq)
		if err != nil {
			fmt.Printf("Error sending DELETE to %s: %v\n", relativePath, err)
			continue
		}
		delResp.Body.Close()

		if delResp.StatusCode == http.StatusNoContent || delResp.StatusCode == http.StatusOK {
			fmt.Printf("✓ Deleted: %s (%s)\n", relativePath,
				map[bool]string{true: "folder", false: "file"}[isDirectory])
		} else {
			fmt.Printf("✗ Failed to delete %s (status %d)\n", relativePath, delResp.StatusCode)
		}
	}

	return nil
}

func PurgePullZoneCache(pullZoneID string) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("%s/pullzone/%s/purgeCache", apiBase, pullZoneID)
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

func CreatePullZone(storageZoneID, name string, eu, us, asia, sa, af bool) (id, defaultHostname string, err error) {
	const apiBase = "https://api.bunny.net"

	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("BUNNYNET_ACCESS_KEY environment variable is not set")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	payload := map[string]any{
		"Name":              name,
		"Type":              0,
		"OriginType":        2,
		"StorageZoneId":     storageZoneID,
		"EnableGeoZoneEU":   eu,
		"EnableGeoZoneUS":   us,
		"EnableGeoZoneASIA": asia,
		"EnableGeoZoneSA":   sa,
		"EnableGeoZoneAF":   af,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiBase+"/pullzone", bytes.NewReader(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("failed to create Pull Zone (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID        int64 `json:"Id"`
		Hostnames []struct {
			Value string `json:"Value"`
		} `json:"Hostnames"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse JSON response: %w\nRaw body: %s", err, string(body))
	}

	if len(result.Hostnames) == 0 {
		return "", "", fmt.Errorf("no hostnames found in response")
	}

	defaultHostname = result.Hostnames[0].Value
	pullZoneID := fmt.Sprintf("%d", result.ID)

	return pullZoneID, defaultHostname, nil
}

func UploadDirectory(zonePassword, storageZoneName, distPath, region string) error {
	var baseURL = "https://" + "se" + ".storage.bunnycdn.com/" + storageZoneName + "/"

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	err := filepath.Walk(distPath, func(localPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(distPath, localPath)
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

func DeletePullZone(pullZoneID string) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("%s/pullzone/%s", apiBase, pullZoneID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting pull zone (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func DeleteStorageZone(storageZoneID string) error {
	const apiBase = "https://api.bunny.net"
	apiKey := os.Getenv("BUNNYNET_ACCESS_KEY")
	if apiKey == "" {
		return fmt.Errorf("BUNNYNET_ACCESS_KEY is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	url := fmt.Sprintf("%s/storagezone/%s", apiBase, storageZoneID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("AccessKey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting storage zone (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func CreateStorageZone(name, region string) (string, string, error) {
	url := "https://api.bunny.net/storagezone"

	payload := map[string]any{
		"Name":            name,
		"Region":          region,
		"StorageZoneType": 0,
		"ZoneTier":        "Standard",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling create payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AccessKey", os.Getenv("BUNNYNET_ACCESS_KEY"))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID       int64  `json:"Id"`
		Password string `json:"Password"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("error parsing JSON response: %w\nRaw body: %s", err, string(body))
	}

	storageZoneID := fmt.Sprintf("%d", result.ID)

	return storageZoneID, result.Password, nil
}
