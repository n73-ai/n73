package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func RmDockerContainer(projectID string) error {
	runCmd := exec.Command("docker", "rm", "-f", projectID)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm -f failed: %s", string(output))
	}
	return nil
}

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
