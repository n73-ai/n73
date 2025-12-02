package fly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CreateApp(projectName string) (string, error) {
  cmd := exec.Command("fly", "apps", "create", projectName)

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

// CreateMachine deploys your Claude app to Fly.io with all required auth files
func CreateMachine(flyConfigPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}

	// Project directory (where your fly.toml lives)
	projectDir := "/home/agust/work/ai/claude"

	// Source paths
	srcClaudeDir := filepath.Join(homeDir, ".claude")
	srcClaudeJSON := filepath.Join(homeDir, ".claude.json")
	srcClaudeBackup := filepath.Join(homeDir, ".claude.json.backup")

	// Destination paths inside the project
	destClaudeDir := filepath.Join(projectDir, ".claude")
	destClaudeJSON := filepath.Join(projectDir, ".claude.json")
	destClaudeBackup := filepath.Join(projectDir, ".claude.json.backup")

	// 1. Clean any previous copies
	_ = os.RemoveAll(destClaudeDir)
	_ = os.Remove(destClaudeJSON)
	_ = os.Remove(destClaudeBackup)

	// 2. Copy the full .claude directory
	if err := copyDir(srcClaudeDir, destClaudeDir); err != nil {
		return fmt.Errorf("failed to copy ~/.claude directory: %w", err)
	}

	// 3. Copy the two JSON files (they must be at the project root)
	if err := copyFile(srcClaudeJSON, destClaudeJSON); err != nil {
		return fmt.Errorf("failed to copy ~/.claude.json: %w", err)
	}
	if err := copyFile(srcClaudeBackup, destClaudeBackup); err != nil {
		// .backup is optional — don't fail if missing
		fmt.Println("Note: ~/.claude.json.backup not found (this is okay)")
	}

	// 4. Ensure sensitive files have correct permissions
	_ = os.Chmod(destClaudeJSON, 0600)
	if _, err := os.Stat(destClaudeBackup); err == nil {
		_ = os.Chmod(destClaudeBackup, 0600)
	}

	// 5. Always clean up everything when done
	defer func() {
		_ = os.RemoveAll(destClaudeDir)
		_ = os.Remove(destClaudeJSON)
		_ = os.Remove(destClaudeBackup)
		fmt.Println("Cleaned up .claude, .claude.json, and .claude.json.backup")
	}()

	// 6. Deploy
	fmt.Println("Deploying to Fly.io with full Claude credentials...")
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

// copyDir recursively copies a directory tree
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

// copyFile copies a single file and preserves permissions where sensible
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	mode := os.FileMode(0644)
	if strings.Contains(src, "credentials") || strings.HasSuffix(src, ".json") {
		mode = 0600
	}
	return os.WriteFile(dst, data, mode)
}
