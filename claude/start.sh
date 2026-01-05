#!/bin/bash

mkdir /root/.claude
cp .credentials.json /root/.claude/.credentials.json

python main.py &

npm --prefix /app/ui-only install
npm --prefix /app/ui-only run dev

