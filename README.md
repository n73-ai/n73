# n73

https://n73.agustfricke.com/verify/

## Build, Preview, and Ship with AI

**n73** is an AI-powered app development platform that lets you create entire web 
applications using **natural language prompts**.

**Version 1** is now available — try it here:
[https://n73.agustfricke.com](https://n73.agustfricke.com)

All pending items (bugs, features, and tasks) are tracked on the issues.

> Questions, feedback, or just interested? Hit me up at my email 
**[hej@agustfricke.com](mailto:hej@agustfricke.com)** your input is always appreciated.

## Tests
```bash
go test -v ./utils -run TestPageExist
```

Ahora estoy usando fly.io para mis backend, ahi se crean las machines
para hacer las apps dentro de n83.dev

el docs de despliege esta en FLY.md

### Admin endpoints(only via api)
check projects
```bash
export JWT=""
curl -i -X GET "https://n73.agustfricke.com/admin/projects" \
     -H "Authorization: Bearer $JWT"
```

check the logs
```bash
curl -i -X GET "http://localhost:8080/logs" | jq
```

prod logs
```bash
curl -i -X GET "https://n73.agustfricke.com/logs" \
  -H "Authorization: Bearer ${JWT}"
```

```bash
curl -i -X GET "https://x73-app.fly.dev/projects/latest"
```
