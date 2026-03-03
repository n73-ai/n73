package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateApp(projectName string) (string, error) {
	cmd := exec.Command(
		"fly",
		"apps",
		"create",
		projectName,
		"--org",
		"x73",
	)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	outBytes, _ := io.ReadAll(stdout)
	errBytes, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("command failed: %w\nstdout: %s\nstderr: %s",
			err, outBytes, errBytes)
	}

	hostname, err := getHostname(projectName)
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func getHostname(appName string) (string, error) {
	cmd := exec.Command("fly", "status", "-a", appName, "--json")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	var resp struct {
		Hostname string `json:"hostname"`
	}

	if err := json.Unmarshal(output, &resp); err != nil {
		return "", fmt.Errorf("failed to parse hostname: %w", err)
	}

	return resp.Hostname, nil
}

func AllocateV6(projectID string) error {
	projectDir := filepath.Join(os.Getenv("ROOT_PATH"), "ts-claude")

	cmd := exec.Command("fly", "ips", "allocate-v6", "--private", "-a", projectID)
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := stdout.String() + stderr.String()
		return fmt.Errorf("fly allocate v6 failed: %w\n--- Output ---\n%s", err, output)
	}

	return nil
}

func CreateMachine(flyConfigPath string) error {
	projectDir := filepath.Join(os.Getenv("ROOT_PATH"), "ts-claude")

	// fly ips allocate-v6 --private

	cmd := exec.Command("flyctl", "deploy", "--config", flyConfigPath, "--ha=false")
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := stdout.String() + stderr.String()
		return fmt.Errorf("fly deploy failed: %w\n--- Output ---\n%s", err, output)
	}

	return nil
}

func GenerateFlyToml(projectID string) error {
	content := fmt.Sprintf(`
app = "%s"
primary_region = 'arn'

[http_service]
  internal_port = 5000
  force_https = true
  auto_stop_machines = 'suspend'
  auto_start_machines = true
  min_machines_running = 0

[[vm]]
  memory = '4gb'
  cpu_kind = 'shared'
  cpus = 4
`, projectID)

	fileName := fmt.Sprintf("%s.toml", projectID)
	flyTomlPath := filepath.Join(os.Getenv("ROOT_PATH"), "fly_configs", fileName)
	return os.WriteFile(flyTomlPath, []byte(content), 0644)
}
