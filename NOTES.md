# notes
just for me!

---

## Deploy!
- add the fly_hostname
```bash
ALTER TABLE projects
ADD COLUMN fly_hostname VARCHAR(255) DEFAULT '';
```
- checkout the new env os.Getenv("DOMAIN") used on CreateProject()
---
