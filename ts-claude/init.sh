#!/bin/sh
set -e

echo "[1/2] Copiando credenciales..."
mkdir -p /root/.claude
cp .credentials.json /root/.claude/.credentials.json

echo "[2/2] Iniciando server..."
npx tsx src/server.ts
