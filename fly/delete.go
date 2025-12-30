package fly

import (
	"fmt"
	"io"
	"os/exec"
)

func DeleteApp(projectID string) error {
	cmd := exec.Command("fly", "apps", "destroy", projectID, "-y")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	outBytes, _ := io.ReadAll(stdout)
	errBytes, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w\nstdout: %s\nstderr: %s",
			err, outBytes, errBytes)
	}

	return nil
}
