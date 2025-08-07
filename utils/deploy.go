package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Push(name, path string) error {
	// push new code to github
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-push.sh")
	cmd := exec.Command(scriptPath, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	// push changes to cloudflare
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	scriptPath = filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-push.sh")
	cmd = exec.Command(scriptPath, slug, path)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}

func FistDeployment(name, path string) error {
	// create and push to remote github repository
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	scriptPath := filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "gh-create-push.sh")
	cmd := exec.Command(scriptPath, slug, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	// create cloudflare page
	slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	scriptPath = filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-create.sh")
	cmd = exec.Command(scriptPath, slug)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	// push the build code to cloudflare
	slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	scriptPath = filepath.Join(os.Getenv("ROOT_PATH"), "scripts", "cf-push.sh")
	cmd = exec.Command(scriptPath, slug, path)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}

func NpmRunBuild(projectPath string) error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	cmd = exec.Command("npm", "run", "build")
	cmd.Dir = projectPath

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}

// !OLD!
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
