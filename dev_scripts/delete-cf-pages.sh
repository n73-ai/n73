#!/bin/bash

# Lista de proyectos a eliminar (todos excepto zustack-video)
projects=(
    "project-96c9ead8-7555-4965-8a87-68ad646fa5da"
    "project-701b6439-91ec-4450-8e21-1bbf63c0bffd"
    "project-164fed94-dd9c-4b47-9233-722ed78423ac"
    "project-d1d61188-b5ec-4270-8983-1e959cae2231"
    "project-3c239f0c-9574-4227-b332-5e1a64ddf45a"
    "project-57916df0-f290-4cba-8489-8ada4d3e7d23"
    "project-66de878a-da29-41c2-805c-4051fa3843fa"
    "project-0e3d2586-db0e-47fd-b0f0-c52445196d66"
    "project-49ad8877-889f-4255-a6d0-336a928b09ed"
)

# Loop para eliminar cada proyecto
for project in "${projects[@]}"; do
    wrangler pages project delete "$project"
done
