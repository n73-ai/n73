#!/bin/bash

# Lista de repositorios a eliminar
repos=(
    "n73-projects/new-tetris-under-n73-projects-list"
)

# Loop para eliminar cada repositorio
for repo in "${repos[@]}"; do
    gh repo delete "$repo" --yes
done
