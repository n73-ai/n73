# AI Zustack

## Build, Preview, and Ship with AI.

⚠️ **This project is under active development. It's not yet production-ready.**

## Simple prompts
- "Build a hello world with a dark background and bold letters"
- "Edit the Hello World text, i want the color to be red"

## todos & tests
- [ ]

# How to use (dev)
- Install claude code
```bash
npm install -g @anthropic-ai/claude-code
```
- Run the db with docker for local development
```bash
docker run --name my-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -d postgres
docker ps -a
docker exec -it my-postgres psql -U postgres -d mydb
```
- Export env
```bash
source .env
```
- Run the go server
```bash
go run cmd/main.go
```

# Dependencies
- gh cli `https://cli.github.com/`
- wrangler cli `https://developers.cloudflare.com/workers/wrangler/install-and-update/`
- Docker (install like in the dotfiles) `https://github.com/agustfricke/dotfiles`
- Go (install like in the dotfiles) `https://github.com/agustfricke/dotfiles`
- Claude code
- python3.10-env

# Deploy
- Install dependencies
- Login gh cli
- Login wrangler 
- Login claude code

- Edit .env file with correct values(db and port)

#### Firewall config
```bash
sudo ufw enable
sudo ufw default deny incoming
sudo ufw default allow outgoing
```
```bash
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
```
```bash
sudo ufw status verbose
```

#### Nginx with certbot
```bash
sudo apt install certbot python3-certbot-nginx nginx -y
certbot --nginx
```

- Remove default config and move the config file for nginx.
```bash
rm /etc/nginx/sites-available/default
rm /etc/nginx/sites-enabled/default
mv ~/ai/server_config/default /etc/nginx/sites-available
sudo ln -s /etc/nginx/sites-available/default /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

- Optional: Check if Nginx is OK
```bash
source .env
go run cmd/main.go
```

#### Systemd 
- Edit the server_config/ai.service with .env and then move the file to systemd
```bash
mv ai/server_config/ai.service /etc/systemd/system
```

- Systemd commands
```bash
service zustack start
service zustack status
service zustack stop
```

#### Create symlinks for applications
sudo ln -s /usr/local/node/bin/node /usr/bin/node
sudo ln -s /usr/local/node/bin/npx /usr/bin/npx

> Questions, feedback, or just interested? Hit me up at ```contact@zustack.com```


