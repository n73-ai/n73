#!/bin/bash

# Lista de proyectos a eliminar (todos excepto zustack-video)
projects=(
    "project-96c9ead8-7555-4965-8a87-68ad646fa5da"
)

# Loop para eliminar cada proyecto
for project in "${projects[@]}"; do
    wrangler pages project delete "$project"
done
