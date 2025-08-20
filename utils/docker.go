package utils

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RefreshCommit() error {
	runCmd := exec.Command("docker", "commit", "claude-server", "base:v1")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker command failed: %s", string(output))
	}

	runCmd = exec.Command("docker", "image", "prune", "-f")
	output, err = runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker command failed: %s", string(output))
	}

	return nil
}

func DockerExists(projectID string) error {
	runCmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", projectID), "--format", "{{.Status}}")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker command failed: %s", string(output))
	}

	status := strings.TrimSpace(string(output))
	if status == "" {
		return fmt.Errorf("container '%s' does not exist", projectID)
	}

	return nil
}

func DockerCloneRepo(projectName, projectID string) error {
	runCmd := exec.Command("docker", "exec", projectID, "rm", "-rf", "/app/project")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker failed: %s", string(output))
	}

	repoToClone := fmt.Sprintf("https://github.com/n73-projects/%s", projectName)
	runCmd = exec.Command("docker", "exec", projectID, "git", "clone", repoToClone, "/app/project")
	output, err = runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker failed: %s", string(output))
	}

	return nil
}

func RmDockerContainer(projectID string) error {
	runCmd := exec.Command("docker", "rm", "-f", projectID)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm -f failed: %s", string(output))
	}
	return nil
}

func PowerOn(projectID string, port int) error {
	runCmd := exec.Command("docker", "start", projectID)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %s", string(output))
	}

	maxAttempts := 30 // máximo 30 intentos (30 segundos)
	for i := 0; i < maxAttempts; i++ {
		// Verificar estado del contenedor
		statusCmd := exec.Command("docker", "inspect", "--format={{.State.Running}}", projectID)
		statusOutput, err := statusCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to check container status: %s", string(statusOutput))
		}

		isRunning := strings.TrimSpace(string(statusOutput)) == "true"
		if !isRunning {
			return fmt.Errorf("container failed to start")
		}

		if IsServiceReady(port) {
			fmt.Println(" ✓ Container is ready!")
			return nil
		}

		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	return nil
}

func CreateDockerContainer(projectID string, port int) error {
	runCmd := exec.Command("docker", "run",
		"-d",
		"--network", "host",
		"-e", fmt.Sprintf("PORT=%v", port),
		"--name", projectID,
		"base:v1")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %s", string(output))
	}

	fmt.Printf("Waiting for container to be ready...")

	// Verificar que el contenedor esté funcionando
	maxAttempts := 30 // máximo 30 intentos (30 segundos)
	for i := 0; i < maxAttempts; i++ {
		// Verificar estado del contenedor
		statusCmd := exec.Command("docker", "inspect", "--format={{.State.Running}}", projectID)
		statusOutput, err := statusCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to check container status: %s", string(statusOutput))
		}

		isRunning := strings.TrimSpace(string(statusOutput)) == "true"
		if !isRunning {
			return fmt.Errorf("container failed to start")
		}

		if IsServiceReady(port) {
			fmt.Println(" ✓ Container is ready!")
			return nil
		}

		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("container did not become ready within %d seconds", maxAttempts)
}

// Función auxiliar para verificar si el servicio está respondiendo
func IsServiceReady(port int) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Health check endpoint
	url := fmt.Sprintf("http://0.0.0.0:%d/health", port)

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

/*
func CreateDockerContainer(projectID string, port int) error {
	ports := fmt.Sprintf("%d:5000", port)
	runCmd := exec.Command("docker", "run", "-d", "-p", ports, "--name", projectID, "base:v1")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %s", string(output))
	}

	fmt.Printf("Waiting for container to initialize...")
	time.Sleep(10 * time.Second)

	return nil
}
*/

func cleanDirectoryExceptGit(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		entryPath := filepath.Join(dirPath, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", entryPath, err)
		}
	}

	return nil
}

func CopyProjectToExisitingProject(projectID string) error {
	destPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)

	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err := cleanDirectoryExceptGit(destPath); err != nil {
		return fmt.Errorf("failed to clean directory: %w", err)
	}

	cmd := exec.Command("docker", "cp", projectID+":/app/project/.", destPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker cp failed: %s", string(output))
	}

	return nil
}

func TryBuildProject(projectID string) error {
	buildCmd := exec.Command("docker", "exec", projectID, "sh", "-c",
		"cd /app/project && npm i && npm run build")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %s", string(output))
	}
	return nil
}
