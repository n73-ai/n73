#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Error usage: $0 <repo-name> <project-path>"
    exit 1
fi

NAME=$1
PROJECT_PATH=$2

cd "$PROJECT_PATH" || exit

if [ ! -d .git ]; then
    git init
    git add .
    git commit -m "Initial commit"
fi

gh repo create "$NAME" --public --source=. --remote=origin --push
