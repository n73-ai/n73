package fly

import (
	"fmt"
	"os/exec"
)

func RebootApp(appName string) error {
	cmd := exec.Command("fly", "apps", "restart", appName)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed exec command fly apps restart <app-name> err: %w. out: %s", err, out)
	}

	return nil
}
