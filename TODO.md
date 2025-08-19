# todo v2
A to-do list of tasks for the v2 of n73 project — feel free to work on these and 
submit a pull request.

### Fix
- [ ] system_prompt: ai olvida de importar cosas a App.tsx(debe hacer los cambios visibles para el usuario)

### Bug
- [ ] why the 401 error? is the script.sh executing on the claude server?
- [ ] when error try to fix loader keeps spinning
- [ ] Prompt is too long(delete docker container when this happends, .claudeignore?)
- [ ] token expired(commit the container more often)
- [ ] resume project: docker not starting because of bad logic !p.DockerRunning is not correct, why?(must delete the docker from the admin, not in the cli)
- [ ] github error:
 error: failed to push some refs to 'https://github.com/n73-projects/project-d849fcf1-0e0f-4e47-b7e2-f078b9ef1099'
 hint: Updates were rejected because the tip of your current branch is behind
 hint: its remote counterpart. Integrate the remote changes (e.g.
 hint: 'git pull ...') before pushing again.
 hint: See the 'Note about fast-forwards' in 'git push --help' for details.
 - [ ] when github error the app stops showing the iframe

### Test
- [ ] how many projects can work at the same time? 

### Feature
- [ ] Delete project in GitHub & Cloudflare
- [ ] User Limit
- [ ] Load balancer
- [ ] Backend integration
- [ ] User profile page
