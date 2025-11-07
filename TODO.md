# todo v1.1.0

A to-do list of tasks for the v1.0.1 of n73 project — feel free to work on these and
submit a pull request.

# done
## Bug Fixes
- Fixed project state handling: the state is now managed only on the server for more accurate control.
- Resolved duplicate spinners when loading project and chat - now only a single spinner is shown.  
- Prevented iframe from refreshing on window focus.  
- Removed missing auth token issue in Claude integration.  
- Remove docker_running from the projects table; check if the service container is available instead
- Update correct ai models 
## Features
- Added log endpoint for admin.  
- Responsive on project.tsx

# Deploy
```sql
ALTER TABLE projects
ADD COLUMN error_msg TEXT DEFAULT '';

ALTER TABLE projects
DROP COLUMN docker_running;
```
- do db migration for the new project.error_msg
- delete docker_running from projects

### Fix
- [ ] 

### Bug
- [ ] 

### Test
- [ ] 

### Feature
- [ ] more fast to deploy the project(go routines)
- [ ] admin panel to power off docker containers
- [ ] User Limit
- [ ] Load balancer
- [ ] Backend integration
- [ ] User profile page
