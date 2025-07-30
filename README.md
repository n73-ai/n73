# AI Zustack

⚠️ **This project is under active development. It's not yet production-ready.**

### Your favorite coding vibe space, powered by AI.

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

## Links and docs
https://docs.anthropic.com/en/docs/claude-code/sdk
https://console.anthropic.com/cost

Questions, feedback, or just interested? Hit me up at ```contact@zustack.com```
