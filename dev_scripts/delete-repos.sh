#!/bin/bash

# Lista de repositorios a eliminar
repos=(
    "n73-projects/project-701b6439-91ec-4450-8e21-1bbf63c0bffd"
    "n73-projects/project-164fed94-dd9c-4b47-9233-722ed78423ac"
    "n73-projects/project-d1d61188-b5ec-4270-8983-1e959cae2231"
    "n73-projects/project-3c239f0c-9574-4227-b332-5e1a64ddf45a"
    "n73-projects/project-57916df0-f290-4cba-8489-8ada4d3e7d23"
    "n73-projects/project-66de878a-da29-41c2-805c-4051fa3843fa"
    "n73-projects/project-0e3d2586-db0e-47fd-b0f0-c52445196d66"
)

# Loop para eliminar cada repositorio
for repo in "${repos[@]}"; do
    gh repo delete "$repo" --yes
done
