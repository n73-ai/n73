# fixes
- [testing] when doing the screenshot, wait a till all the content is loaded(it's not loading the images).
- [ ] system prompt: que no diga, corre este comando o ejecuta esto, ya que el usuario no lo puede hacer.

# feats
- [ ] upload projects to github.
    - [ ] is main server available to resive source code?
        - check:
            disk_available = disk size - disk usage
            if disk_available <(mayor que):
                download & commit & push
            else:
                send_from_remote(project.gh_status = pending)
        - go routine:
            every hour check if there is some project with project_id.gh_status == pending)
            if projects.length < 0:
                for p in projects:
                    send request to power on remote server and run check
    - [ ] send source code to main server
    - [ ] script to commit & push to github
- [ ] add custom domain
- [ ] upload images to claude-code
- [ ] generate buffer of fly apps to generate the projects faster.
- [ ] supabase 
- [ ] list of commits history from ui
- [ ] rollback to commit
- [ ] visual edits
