package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CreateCFPage(name string) error {
  slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-create.sh")
	cmd := exec.Command(scriptPath, slug)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func CreatePushGH(name, path string) error {
  slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
  fmt.Println("th slug", slug)
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-create-push.sh")
	cmd := exec.Command(scriptPath, slug, path)
  fmt.Println("the command", cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func PushCF(name, path string) error {
  slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-push.sh")
	cmd := exec.Command(scriptPath, slug, path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func PushGH(path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-push.sh")
	cmd := exec.Command(scriptPath, path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
