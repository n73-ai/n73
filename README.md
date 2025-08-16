# n73

## Build, Preview, and Ship with AI

**n73** is an AI-powered app development platform that lets you create entire web applications using **natural language prompts**.
Transform your ideas into reality with the power of artificial intelligence.

**Version 1** is now available — try it here:
[https://n73.agustfricke.com](https://n73.agustfricke.com)

All pending items (bugs, features, and tasks) are tracked in the TODO.md file.

> Questions, feedback, or just interested? Hit me up at my email **[hej@agustfricke.com](mailto:hej@agustfricke.com)** your input is always appreciated.

---

## Local Setup (Ubuntu 22.04)

### Dependencies

#### Node.js

```bash
wget https://nodejs.org/dist/v22.17.1/node-v22.17.1-linux-x64.tar.xz
tar -xf node-v22.17.1-linux-x64.tar.xz
rm node-v22.17.1-linux-x64.tar.xz
mv node-v22.17.1-linux-x64 /usr/bin/node
```

Add to your `.bashrc`:

```bash
export PATH=$PATH:/usr/bin/node/bin
```

Apply changes:

```bash
source /path/to/.bashrc
```

#### Wrangler

```bash
npm i -D -g wrangler@latest
```

#### GitHub CLI

```bash
(type -p wget >/dev/null || (sudo apt update && sudo apt install wget -y)) \
 && sudo mkdir -p -m 755 /etc/apt/keyrings \
 && out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
 && sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg < $out > /dev/null \
 && sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
 && sudo mkdir -p -m 755 /etc/apt/sources.list.d \
 && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
    | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
 && sudo apt update \
 && sudo apt install gh -y
```

#### Docker

```bash
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
 | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
 https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" \
 | sudo tee /etc/apt/sources.list.d/docker.list
sudo apt-get update -y
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
```

#### Go

```bash
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
tar -xf go1.22.5.linux-amd64.tar.gz
rm go1.22.5.linux-amd64.tar.gz
mv go /usr/bin
```

Add to your `.bashrc`:

```bash
export PATH=/usr/bin/go/bin:$PATH
export GOPATH=/.go
export PATH=$PATH:$GOPATH/bin
```

Apply changes:

```bash
source /path/to/.bashrc
```

#### Python 3.10 Environment

```bash
apt install python3.10-venv
```

---

### Authentication & Login Setup

#### Generate SSH Key

```bash
ssh-keygen -t ed25519 -C "hej@agustfricke.com"
```

Start ssh-agent and add key:

```bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
```

Copy your public key:

```bash
cat ~/.ssh/id_ed25519.pub
```

Add it to your GitHub account.

#### GitHub CLI Login

```bash
gh auth login
```

#### Wrangler Login

```bash
wrangler login
```

Follow the OAuth flow — after it redirects to `localhost:[port]`,
open a new terminal on your server and run:

```bash
curl "localhost_url"
```

---

### Claude Code (Docker Build & Commit)

```bash
docker build -t claude-server .
docker run --network host -d  --name claude-server claude-server
# new
docker run --network host -d  --name claude-server claude-server

# new with port, dont need the port because if not port, sets the port to 5000
docker run --network host -e PORT=5000 -d --name claude-server claude-server
# end new

docker exec -it claude-server bash
claude  # Login to Claude
docker commit claude-server base:v1
```

---

### Database Setup

psql -h 178.18.240.78 -p 5432 -U postgres -d mydb

```bash
docker run -d \
  --name owo \
  --network host \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=mydb \
  postgres

docker run --name n73-database \
  --network host \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=mydb \
  -d postgres

  --volumes-from my-postgres \

  -p 5432:5432 \
docker exec -it n73-database psql -U postgres -d mydb

psql -h localhost -p 5432 -U postgres -d mydb

docker run --name my-new-postgres \
  --network host \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=mydb \
  -e PGPORT=5433 \
  -d postgres
```

Paste the contents of `tables.sql` into the PostgreSQL shell.

---

### Environment Variables

**Frontend** (`ai/ui/.env.local`):

```bash
VITE_BACKEND_URL="http://localhost:8080"
VITE_WEBSOCKET_URL="ws://localhost:8080"
```

**Backend**:

```bash
export DB_USER=postgres
export DB_PASSWORD=secret
export DB_HOST=localhost
export DB_NAME=mydb
export DB_PORT=5432
export EMAIL_SECRET_KEY=secret_email_pass
export PORT=8080
export ROOT_PATH=/path/to/ai
export SECRET_KEY="very_long_string"
export ADMIN_JWT="jwt_with_email_admin_user_claim"
export IP="192.168.1.9"
```

---

### Notes

* **EMAIL\_SECRET\_KEY**: Obtain from your Google account to enable email sending.
* **ROOT\_PATH**: Path to your AI project (e.g., `/home/agust/work/ai`).
* **ADMIN\_JWT**: Generate using this Go snippet:

```go
package main

import (
  "fmt"
  "time"

  "github.com/golang-jwt/jwt"
)

func main() {
  tokenByte := jwt.New(jwt.SigningMethodHS256)
  now := time.Now().UTC()
  claims := tokenByte.Claims.(jwt.MapClaims)
  //expDuration := time.Hour * 24 * 180
  //exp := now.Add(expDuration).Unix()
  //claims["exp"] = exp
  claims["email"] = "agustfricke@gmail.com"
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()
  tokenString, err := tokenByte.SignedString([]byte("sodibg3obg48wb4ogbwsl4gjwnf4owbg9snkfbitbfwouebgiw893r83bf9uw3bsfeiugfiwbf4ifgw938hg9w74gfi4wubcwih39h0f038298yhw8beguiwebgiwbe=w-eto293ru2094tf-w=efwnoigb"))
  if err != nil {
    panic(err)
  }
  fmt.Println(tokenString)
}
```

Run it and use the output as your admin JWT.

* **IP**: Your local IP (use `hostname -I`).

---

### Running the Project

**Backend:**

```bash
go run cmd/main.go
```

Create a new user in the UI, then make them admin:

```bash
docker exec -it my-postgres psql -U postgres -d mydb
```

```sql
UPDATE users SET role = 'admin' WHERE email = 'your@admin.com';
```

**Frontend:**

```bash
cd ui
npm run dev
```

Visit [http://localhost:5173](http://localhost:5173) to use the app locally.

> You’ll also need to update the repository upload location.
> In my case, it’s under **GitHub `n73-projects`**, so make sure to change this in the `scripts/gh-create.sh` file.

## Deployments

### 1. Firewall Configuration

Start a `screen` session to keep the firewall setup running even if you disconnect:

```bash
screen -S firewall-setup
```

Configure **UFW**:

```bash
sudo ufw allow ssh
sudo ufw enable
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 80
sudo ufw allow 443
sudo ufw status verbose
```

---

### 2. Nginx + Certbot (HTTPS & WSS)

Point your DNS record to your server's IP address, then install Nginx and Certbot:

```bash
sudo apt install certbot python3-certbot-nginx nginx -y
certbot --nginx
```

Replace the default Nginx config with your custom config:

```bash
rm /etc/nginx/sites-available/default
rm /etc/nginx/sites-enabled/default
cp ~/ai/server_config/default /etc/nginx/sites-available
sudo ln -s /etc/nginx/sites-available/default /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

**Optional:** Test if Nginx and the application work as expected:

```bash
cd ~/ai

# Build frontend
cd ui
npm run build

# Run backend
cd ..
source .env
go run cmd/main.go
```

---

### 3. Systemd Service

**Step 1:** Edit `server_config/ai.service` and update environment variables (`.env`).
**Step 2:** Move the file to systemd:

```bash
mv ai/server_config/ai.service /etc/systemd/system
```

**Step 3:** Manage the service:

```bash
sudo service ai start
sudo service ai status
sudo service ai stop
```
