#!/bin/bash

if [ -z "$1" ]; then
  echo "Error ./cf-create.sh <name>"
  exit 1
fi

NAME=$1

wrangler pages project create $NAME --production-branch "main" 
