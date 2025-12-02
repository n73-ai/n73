#!/bin/bash

mv .claude /root
mv .claude.json /root
mv .claude.json.backup /root

python keep_alive.py &

python main.py
