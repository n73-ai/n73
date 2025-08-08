# AI Zustack

## Build, Preview, and Ship with AI.

⚠️ **This project is under active development. It's not yet production-ready.**

## Simple prompts
- "Build a hello world with a dark background and bold letters"
- "Edit the Hello World text, i want the color to be red"

## todos 
- [ ] test in remote server via 8080
- [ ] systemd for python server?
- [ ] if error at resume, create project or webhook set correct state for the ui.
- [ ] container for ai projects
- [ ] better docs

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
- Config .env.local on react
```bash
vim ui/.env.local
```
- Export env
```bash
source .env
```
- Run the go server
```bash
go run cmd/main.go
```

export the env for python server
```bash
export ADMIN_JWT_TOKEN=sign_admin_token
```
```bash
python3 main.py
```

- Make admin user
```bash
UPDATE users SET role = 'admin' WHERE email = 'agustfricke@gmail.com';
```

# Dependencies
- gh cli `https://cli.github.com/`
- wrangler cli `https://developers.cloudflare.com/workers/wrangler/install-and-update/`
- Docker (install like in the dotfiles) `https://github.com/agustfricke/dotfiles`
- Go (install like in the dotfiles) `https://github.com/agustfricke/dotfiles`
- Claude code
- python3.10-env

# Deploy
#### Install dependencies
For `Ubuntu 22.04`.
- Go
```bash
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
tar -xf go1.22.5.linux-amd64.tar.gz
rm go1.22.5.linux-amd64.tar.gz
mv go /usr/bin
```
- Node
```bash
wget https://nodejs.org/dist/v22.17.1/node-v22.17.1-linux-x64.tar.xz
tar -xf node-v22.17.1-linux-x64.tar.xz
rm node-v22.17.1-linux-x64.tar.xz
mv node-v22.17.1-linux-x64 /usr/bin/node
```
Add to .bashrc
```bash
export PATH=/usr/bin/go/bin:$PATH
export GOPATH=/.go
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:/usr/bin/node/bin
```
- Docker
```bash
apt-get remove docker docker-engine docker.io containerd runc 
apt-get update 
apt-get install ca-certificates curl gnupg 
install -m 0755 -d /etc/apt/keyrings 
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg  
chmod a+r /etc/apt/keyrings/docker.gpg 
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" | sudo tee /etc/apt/sources.list.d/docker.list 
apt-get update -y 
apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y 
```
- Wrangler
```bash
npm i -D -g wrangler@latest
```
gh cli
```bash
(type -p wget >/dev/null || (sudo apt update && sudo apt install wget -y)) \
	&& sudo mkdir -p -m 755 /etc/apt/keyrings \
	&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
	&& cat $out | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
	&& sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
	&& sudo mkdir -p -m 755 /etc/apt/sources.list.d \
	&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
	&& sudo apt update \
	&& sudo apt install gh -y
```
- Python venv
```bash
# check python version
python3 --version
# if python 3.12.* set python3.12-venv
apt install python3.12-venv
```
- Claude code
```bash
npm install -g @anthropic-ai/claude-code
```
- check if install is correct
```bash
node --version
go version
docker --version
wrangler --version
gh --version
```

#### Auth
##### GitHub ssh
Generating a new SSH key
```bash
ssh-keygen -t ed25519 -C "hej@agustfricke.com"
```
Start the ssh-agent in the background and add your SSH private key to the ssh-agent
```bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
```
Copy and past to your GitHub account
```bash
cat ~/.ssh/id_ed25519.pub
```
##### gh 
```bash
gh auth login
```

##### wrangler
```bash
```
Go to your browser and open the oauth url
it will redirect to a localhost:[port]
now, open a new terminal in your server and curl that url
```bash
curl "localhost_url"
```

##### claude code
```bash
claude
```
Open the url in local web browser, copy the code and pase in the terminal

#### Clone the repo
```bash
git clone git@github.com:zustack/ai.git
```

- Edit .env file with correct values(db and port)

- Config .env.local on react
```bash
vim ui/.env.local
```
- Config the postgres db and add credentials to .env

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


