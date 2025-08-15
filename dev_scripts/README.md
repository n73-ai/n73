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

│ Project Name                                 │ Project Domains                                        │ Git Provider │ Last Modified  │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-701b6439-91ec-4450-8e21-1bbf63c0bffd │ project-701b6439-91ec-4450-8e21-1bbf63c0bffd.pages.dev │ No           │ 24 seconds ago │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-164fed94-dd9c-4b47-9233-722ed78423ac │ project-164fed94-dd9c-4b47-9233-722ed78423ac.pages.dev │ No           │ 1 hour ago     │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-d1d61188-b5ec-4270-8983-1e959cae2231 │ project-d1d61188-b5ec-4270-8983-1e959cae2231.pages.dev │ No           │ 1 hour ago     │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-3c239f0c-9574-4227-b332-5e1a64ddf45a │ project-3c239f0c-9574-4227-b332-5e1a64ddf45a.pages.dev │ No           │ 2 hours ago    │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-57916df0-f290-4cba-8489-8ada4d3e7d23 │ project-57916df0-f290-4cba-8489-8ada4d3e7d23.pages.dev │ No           │ 2 hours ago    │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-66de878a-da29-41c2-805c-4051fa3843fa │ project-66de878a-da29-41c2-805c-4051fa3843fa.pages.dev │ No           │ 2 hours ago    │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-0e3d2586-db0e-47fd-b0f0-c52445196d66 │ project-0e3d2586-db0e-47fd-b0f0-c52445196d66.pages.dev │ No           │ 3 hours ago    │
├──────────────────────────────────────────────┼────────────────────────────────────────────────────────┼──────────────┼────────────────┤
│ project-49ad8877-889f-4255-a6d0-336a928b09ed │ project-49ad8877-889f-4255-a6d0-336a928b09ed.pages.dev │ No           │ 3 hours ago    │

n73-projects/project-701b6439-91ec-4450...               public  less than a minute ago
n73-projects/project-164fed94-dd9c-4b47...               public  about 1 hour ago
n73-projects/project-d1d61188-b5ec-4270...               public  about 1 hour ago
n73-projects/project-3c239f0c-9574-4227...               public  about 2 hours ago
n73-projects/project-57916df0-f290-4cba...               public  about 2 hours ago
n73-projects/project-66de878a-da29-41c2...               public  about 2 hours ago
n73-projects/project-0e3d2586-db0e-47fd... 
