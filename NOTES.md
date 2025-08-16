# Just notes of code and debug here.

docker exec -it claude-server bash

### run docker in same network
docker run --network host -d  --name claude-server claude-server

- is the docker running?
curl -I http://localhost:5000/health

33137

- Delete the commit if needit
```bash
docker rmi base:v1
```


## fixing the firewall bug

edits: docker.go, main.py

pasar el puerto del servior de python a travez de una env
```go
runCmd := exec.Command("docker", "run",
    "-d",
    "--network", "host",
    // NO uses "-p" con network host
    "-e", fmt.Sprintf("PORT=%s", port), // Pasar puerto como variable de entorno
    "--name", projectID,
    "base:v1")
output, err := runCmd.CombinedOutput()
if err != nil {
    return fmt.Errorf("docker run failed: %s", string(output))
}
```

cambiar esto en el codigo de python, poner el puerto desde env
```py
import os

# Leer puerto de variable de entorno, por defecto 5000
port = int(os.getenv('PORT', 5000))

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=port, debug=True)
```

## 401: OAuth token has expired under claude
func() {
    `docker commit claude-server base:v1` volver a hacer el commit de la misma imagen para que el token este fresssco
    `docker image prune` para elminar todos los commits de claude-server y ->
}()
