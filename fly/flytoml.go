package fly

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateFlyToml(projectID string) error {
	content := fmt.Sprintf(`
app = "%s"
primary_region = 'arn'

[http_service]
  internal_port = 5173
  force_https = true
  auto_stop_machines = 'stop'
  auto_start_machines = true
  min_machines_running = 0

[[vm]]
  memory = '2gb'
  cpu_kind = 'shared'
  cpus = 1
`, projectID)

	fileName := fmt.Sprintf("%s.toml", projectID)
	flyTomlPath := filepath.Join(os.Getenv("ROOT_PATH"), "fly_configs", fileName)
	return os.WriteFile(flyTomlPath, []byte(content), 0644)
}
