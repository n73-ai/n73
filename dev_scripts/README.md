# scripts to bulk delete cf pages & gh repos

list
```bash
wrangler pages project list
gh repo list n73-projects
```

then add the projects in bash script to delete then in a for loop

```bash
./delete*
```
