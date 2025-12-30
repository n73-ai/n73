#!/bin/sh
set -e

# Función para descomprimir un tar base64 si existe el secret
extract_tar() {
    local secret_name=$1
    local target_dir=$2

    if [ -n "$(eval echo "\$${secret_name}")" ]; then
        echo "Descomprimiendo $secret_name en $target_dir..."
        mkdir -p "$target_dir"
        # El secret viene en base64 → decodificar → pipe a tar
        eval echo "\$${secret_name}" | base64 -d | tar -xz -C "$target_dir" --strip-components=1
    fi
}

# Descomprimir directorios completos desde secrets (si existen)
extract_tar WRANGLER_TAR /root/.config/wrangler
extract_tar GH_TAR /root/.config/gh
extract_tar CLAUDE_TAR /root/.claude

# Archivos individuales (si prefieres algunos así)
if [ -n "$SSH_PRIVATE_KEY" ]; then
    mkdir -p /root/.ssh
    echo "$SSH_PRIVATE_KEY" > /root/.ssh/id_ed25519
    chmod 600 /root/.ssh/id_ed25519
fi

if [ -n "$CLAUDE_JSON" ]; then
    echo "$CLAUDE_JSON" > /root/.claude.json
fi

if [ -n "$CLAUDE_JSON_BACKUP" ]; then
    echo "$CLAUDE_JSON_BACKUP" > /root/.claude.json.backup
fi

# Ejecutar la aplicación
exec "$@"
