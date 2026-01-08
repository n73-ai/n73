package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DeleteGhRepo(projectID string) error {
	query := fmt.Sprintf("n73-projects/project-%s", projectID)
	cmd := exec.Command("gh", "repo", "view", query, "--json", "name")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "GraphQL: Could not resolve to a Repository with the name") {
			return nil
		}
		return fmt.Errorf("e1: %s", string(output))
	}

	cmd = exec.Command("gh", "repo", "delete", query, "--yes")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("e2%s", string(output))
	}

	return nil
}

func DeleteCfPage(projectID string) error {
	cmd := exec.Command("wrangler", "pages", "project", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	query := fmt.Sprintf("project-%s", projectID)
	if !strings.Contains(string(output), query) {
		return nil
	}

	cmd = exec.Command("wrangler", "pages", "project", "delete", query, "--yes")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}

func GhCreate(slug, projectPath string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-create.sh")
	cmd := exec.Command(scriptPath, slug, projectPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func GhPush(path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-push.sh")
	cmd := exec.Command(scriptPath, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func PageExists(projectName string) (bool, error) {
	cmd := exec.Command("wrangler", "pages", "project", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("error exec wrangler: %s", out.String())
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, projectName) {
			return true, nil
		}
	}

	return false, nil
}

func CfCreate(slug string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-create.sh")
	cmd := exec.Command(scriptPath, slug)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func NpmRunBuild(path string) error {
  installCmd := exec.Command("npm", "install")
  installCmd.Dir = path
  if err := installCmd.Run(); err != nil {
    return fmt.Errorf("npm install failed: %w", err)
  }

  buildCmd := exec.Command("npm", "run", "build")
  buildCmd.Dir = path
  if err := buildCmd.Run(); err != nil {
      return fmt.Errorf("npm run build failed: %w", err)
  }

  return nil
}

func CfPush(slug, path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-push.sh")
	cmd := exec.Command(scriptPath, slug, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
    fmt.Println(err.Error())
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func GhClone(repo, path, projectID string) error {
  projectPath := filepath.Join(path, projectID)
	if _, err := os.Stat(projectPath); err == nil {
		if err := os.RemoveAll(projectPath); err != nil {
			return fmt.Errorf("no se pudo eliminar el directorio existente: %w", err)
		}
	}

	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-clone.sh")
	cmd := exec.Command(scriptPath, repo, path, projectID)
  fmt.Println("gh clone cmd: ", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
