#!/bin/bash

if [ -z "$1" ]; then
  echo "Error ./cf-create.sh <name>"
  exit 1
fi

PROJECT_PATH=$1

git -C $PROJECT_PATH add $PROJECT_PATH
git -C $PROJECT_PATH commit -m "commit from n73.io"
git -C $PROJECT_PATH push 
