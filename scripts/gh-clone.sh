#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <repo-name> <project-path> <projectID>"
    exit 1
fi

REPO=$1          
PROJECT_PATH=$2       
PROJECT_ID=$3       

if [ ! -d "$PROJECT_PATH" ]; then
    exit 1
fi

cd "$PROJECT_PATH" || exit 1

git clone "$REPO" "$PROJECT_ID"
