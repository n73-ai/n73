#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Error usage: $0 <input-1> <input-2>"
    exit 1
fi

NAME=$1
PROJECT_PATH=$2 # ./dist

cd $PROJECT_PATH
npm i
npm run build 
wrangler pages deploy "$PROJECT_PATH/dist" --project-name=$NAME --branch=main
