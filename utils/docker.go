package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func CreateDockerContainer(projectID string) (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	addr := listener.Addr().(*net.TCPAddr)
	assignedPort := addr.Port
	listener.Close() 

	ports := fmt.Sprintf("%d:5000", assignedPort) 
	runCmd := exec.Command("docker", "run", "-d",
		"-p", ports, "--name", projectID, "base:v1")
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("docker run failed: %s", string(output))
	}

	fmt.Printf("Waiting for container to initialize...")
	time.Sleep(10 * time.Second)

	return assignedPort, nil
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

func CopyProjectToMainMachine(projectID string) error {
	destPath := filepath.Join(os.Getenv("ROOT_PATH"), "projects", projectID)

	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
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

func IsContainerRunning(projectID string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "--format={{.State.Running}}", projectID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("container inspection failed: %s", string(output))
	}
	return string(output) == "true\n", nil
}

func SafeCopyAndBuild(projectID string) error {
	running, err := IsContainerRunning(projectID)
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}
	if !running {
		return fmt.Errorf("container %s is not running", projectID)
	}

	if err := CopyProjectToMainMachine(projectID); err != nil {
		return fmt.Errorf("failed to copy project: %w", err)
	}

	if err := TryBuildProject(projectID); err != nil {
		return fmt.Errorf("failed to build project: %w", err)
	}

	return nil
}
