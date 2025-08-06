package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func CreateCFPage(name string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "cf-create.sh")
	cmd := exec.Command(scriptPath, name)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func CreatePushGH(name, path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "gh-create-push.sh")
	cmd := exec.Command(scriptPath, name, path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func PushCF(name, path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "cf-push.sh")
	cmd := exec.Command(scriptPath, name, path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func PushGH(path string) error {
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "gh-push.sh")
	cmd := exec.Command(scriptPath, path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
