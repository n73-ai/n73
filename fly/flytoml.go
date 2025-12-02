package fly

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateFlyToml crea un fly.toml perfecto con el nombre de app que tú elijas
func GenerateFlyToml(projectID string) error {
    // Limpiamos el nombre (Fly solo permite minúsculas, números y guiones)
    content := fmt.Sprintf(`
app = "%s"
primary_region = "cdg"

[build]
  image = "" 

[http_service]
  internal_port = 5000
  force_https = true
  auto_stop_machines = "suspend"
  auto_start_machines = true
  min_machines_running = 0

[[vm]]
  memory = "1gb"
  cpu_kind = "shared"
  cpus = 1

[mounts]
  source = "data"
  destination = "/persist"
`, projectID)

    fileName := fmt.Sprintf("%s.toml", projectID)
    flyTomlPath := filepath.Join(os.Getenv("ROOT_PATH"), "fly_configs", fileName)
    return os.WriteFile(flyTomlPath, []byte(content), 0644)
}
