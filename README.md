# AI Zustack

## Build, Preview, and Ship with AI.

Create a todo website, where users can add, remove and update the task state (done, or no done), make the app in theme dark .

⚠️ **This project is under active development. It's not yet production-ready.**

## Simple prompts
- "Build a hello world with a dark background and bold letters"
- "Edit the Hello World text, i want the color to be red"

## todos & tests
- [ ] add better system prompt
- [ ] imprube basic ui
- [ ] CI/CD docs

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

## Project status:
Building(ai output ongoing)
Ready(finish ai output, npm run build success so is ready to deploy)
Error
Deployed(successfull deployment)
Deploying(doing the deploy)

## system prompt edit: 
- build solo con react, typescript, tailwindcss y shadcn
- puedes instalar componentes de shadcn si crees que es apropiado
- solo menciona lo que estas haciendo dentro del projecto
- solo edita los archivos dentro del react project
- no corras npm run dev 
- 

> Questions, feedback, or just interested? Hit me up at ```contact@zustack.com```
