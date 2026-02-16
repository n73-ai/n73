#!/bin/bash

mkdir /root/.claude
cp .credentials.json /root/.claude/.credentials.json

python main.py &

npm --prefix /app/ui-only install
npm run --prefix /app/ui-only build
python3 -m http.server 5173 --directory /app/ui-only/dist
