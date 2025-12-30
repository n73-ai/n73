#!/bin/bash

mkdir /root/.claude
cp .credentials.json /root/.claude/.credentials.json

python main.py &

npm --prefix /app/project install
npm --prefix /app/project run dev

