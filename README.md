# AI Zustack

## Build, Preview, and Ship with AI.

⚠️ **This project is under active development. It's not yet production-ready.**

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

# Links and docs
(https://docs.anthropic.com/en/docs/claude-code/sdk)
(https://docs.anthropic.com/en/docs/about-claude/models/overview)
(https://console.anthropic.com/cost)

# Dependencies
- gh cli `https://cli.github.com/`
- wrangler cli `https://developers.cloudflare.com/workers/wrangler/install-and-update/`
- Docker (install like in the dotfiles) `https://github.com/agustfricke/dotfiles`
- Go (install like in the dotfiles) `https://github.com/agustfricke/dotfiles`

> Questions, feedback, or just interested? Hit me up at ```contact@zustack.com```
