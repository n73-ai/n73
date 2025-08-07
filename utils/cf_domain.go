package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// domain, err := utils.GetProjectDomainFallback("project-name")
func GetProjectDomainFallback(projectName string) (string, error) {
	cmd := exec.Command("wrangler", "pages", "project", "list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, projectName) && strings.Contains(line, "│") {
			parts := strings.Split(line, "│")
			if len(parts) >= 3 {
				domains := strings.TrimSpace(parts[2])
				domainList := strings.Split(domains, ",")
				for _, d := range domainList {
					d = strings.TrimSpace(d)
					if !strings.HasSuffix(d, ".pages.dev") && d != "" {
						return "https://" + d, nil
					}
				}
				for _, d := range domainList {
					d = strings.TrimSpace(d)
					if strings.HasSuffix(d, ".pages.dev") {
						return "https://" + d, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("project not found")
}
