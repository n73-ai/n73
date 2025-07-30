# AI Zustack

## Build, Preview, and Ship with AI.

⚠️ **This project is under active development. It's not yet production-ready.**

# How to use
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

# todo
- [ ] new project
- [ ] resume project
- [ ] list messages nice styles
    - [ ] thinking...
    - [ ] readme
    - [ ] Zustack {response}
- [ ] list projects on /
- [ ] better system_prompt
- [ ] make base tamplate(react, typescript, tailwindcss)
- [ ] npm run build -> if error button(try to fix)
    - [ ] new 
    - [ ] resume 
- [ ] push to github
    - [ ] new 
    - [ ] resume 
- [ ] deploy to cloudflare
- [ ] when deploy update domain and show on preview
    - [ ] new 
    - [ ] resume 
- [ ] login
- [ ] signup(close on prod)
- [ ] deploy beta
    - [ ] ui cloudflare
    - [ ] backend cloudflare tunnel

# Links and docs
```
https://docs.anthropic.com/en/docs/claude-code/sdk
```
```
https://console.anthropic.com/cost
```

> Questions, feedback, or just interested? Hit me up at ```contact@zustack.com```
