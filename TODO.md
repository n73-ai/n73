# todo v2
A to-do list of tasks for the v2 of n73 project — feel free to work on these and 
submit a pull request.

### Fix
- [ ] change middleware.User for middleware.Admin in require routes

### Bug
- [x] claude server is not responding because of timeout with host ip(try to change the ip for the domain name nginx proxy)
- [ ] firewall not working(can access docker containers from remote)
- [ ] 401: OAuth token has expired under claude(recommit)
- [ ] when not logged in must redirect to login page
- [ ] Spinner in client keeps spinning when finish(1.error try to fix, 2.page deployed)
- [ ] Remote machine out of memory: "--memory=300m" & "--cpus=0.5" when creating a new docker container
    - [ ] how many machines can work at the same time? based on that i can build the load balancer

### Feature
- [ ] system prompt: no te olvides de importar a App.tsx
- [ ] User Limit
- [ ] Load balancer
- [ ] Backend integration
- [ ] User profile page
