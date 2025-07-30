package utils

import (
	"fmt"
	"os"
)

func CheckRequiredEnv(vars []string) error {
	var missing []string

	for _, v := range vars {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v.", missing)
	}

	return nil
}
