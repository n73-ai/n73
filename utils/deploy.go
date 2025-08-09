package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
