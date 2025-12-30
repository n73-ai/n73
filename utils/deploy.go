package utils

import (
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

/*
// remove this function and test if it works
func DeleteRemote(projectID string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "delete-remote.sh")
	cmd := exec.Command(scriptPath, projectID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
*/

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

func CfCreate(slug string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-create.sh")
	cmd := exec.Command(scriptPath, slug)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func CfPush(slug, path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-push.sh")
	cmd := exec.Command(scriptPath, slug, path)
	fmt.Println("cmd: ", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
