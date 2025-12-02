# todo v1.1.0

A to-do list of tasks for the next version of n73 — feel free to work on these and
submit a pull request.

### Features
- [x] auth claude in machine
- [x] Run fly machine
- [x] compress the project directory, then send it over the network to the main server
      push to github and then push to cloudflare, use it in messages
      cuando termine de escribir el codigo, manda el .zip y lo recibe en messages
      ahi lo decomprime y hace push(para tener el estado el projecto de manera correcta)
- [x] python server in fly machine
- [x] resume project
- [x] update claude models
- [x] delete project(delete also from fly.io)
- [x] better system prompt

- [ ] transfer project to other user
- [ ] stripe payments

- [ ] migrate code from react to next.js like in the sample of ~/personal/nextjs

## branches

## preformace
- [ ] Add docker build to registry for fast machines
- [ ] Create buffer of machines for fast code generation at first prompt(give claude some context)

## backend feat
- [ ] backend go base codebase 
- [ ] Postgres(fly io)
- [ ] Tigris(fly io)
- [ ] poweroff claude server and poweron go server in dev mode

[dev_app1] claude-code code generation | backend dev server | Postgres
[prod_app2] backend production server | Postgres

cuando la ai termino de responder,
mandar el codigo?

go test -v ./fly
go test -v ./utils -run TestUnzip

### Fix
- [ ] fix vulnerabilities listed on "https://www.shodan.io/host/178.18.240.78" 

### Bugs
- [ ] 

### Tests
- [ ] 
