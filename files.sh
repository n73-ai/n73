#!/bin/bash
set -e

# Función para copiar al portapapeles (cross-platform)
copy_to_clipboard() {
    if command -v pbcopy >/dev/null; then
        pbcopy
    elif command -v xclip >/dev/null; then
        xclip -selection clipboard
    elif command -v wl-copy >/dev/null; then
        wl-copy
    else
        echo "⚠️  No se detectó portapapeles. Copia manualmente el texto de arriba."
        cat
        return 1
    fi
}

# Función para pegar del portapapeles
paste_from_clipboard() {
    if command -v pbpaste >/dev/null; then
        pbpaste
    elif command -v xclip >/dev/null; then
        xclip -selection clipboard -o
    elif command -v wl-paste >/dev/null; then
        wl-paste
    else
        echo "⚠️  No se detectó portapapeles. Pega manualmente el contenido."
        return 1
    fi
}

# Función para subir un secret empaquetado en tar.gz + base64
upload_tar_secret() {
    local name="$1"
    local path="$2"
    local description="$3"

    echo "📦 Generando $name desde $path..."

    if [ ! -e "$path" ]; then
        echo "❌ No existe: $path → Saltando"
        return
    fi

    # Crear tar.gz en memoria y codificar en base64
    local content=$(tar -cz -C "$(dirname "$path")" "$(basename "$path")" 2>/dev/null | base64 -w 0)

    if [ -z "$content" ]; then
        echo "❌ El directorio está vacío o falló el tar: $path"
        return
    fi

    echo "✅ $description generado ($(( ${#content} / 1024 )) KB)"
    echo "$content" | copy_to_clipboard && echo "📋 Copiado al portapapeles!"

    read -p "¿Quieres subirlo ahora como secret Fly.io? (y/N): " confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        echo "$content" | fly secrets set "$name"=- || echo "❌ Falló fly secrets set"
    fi
    echo
}

# Función para archivos individuales
upload_file_secret() {
    local name="$1"
    local path="$2"
    local description="$3"

    echo "📄 Procesando $description..."

    if [ ! -f "$path" ]; then
        echo "❌ No existe: $path → Saltando"
        return
    fi

    local content=$(cat "$path" | base64 -w 0)
    echo "✅ $description generado"
    echo "$content" | copy_to_clipboard && echo "📋 Copiado al portapapeles!"

    read -p "¿Subir como secret $name? (y/N): " confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        echo "$content" | fly secrets set "$name"=- || echo "❌ Error"
    fi
    echo
}

# === Aquí pones lo que realmente necesitas subir ===

upload_tar_secret "WRANGLER_TAR" "$HOME/.config/.wrangler" "config de Wrangler (directorio completo)"
upload_tar_secret "GH_TAR"       "$HOME/.config/gh"      "config de GitHub CLI (gh)"
upload_tar_secret "CLAUDE_TAR"   "$HOME/.claude"         "directorio .claude (si existe)"

upload_file_secret "SSH_PRIVATE_KEY"      "$HOME/.ssh/id_ed25519"       "clave SSH privada"
upload_file_secret "CLAUDE_JSON"          "$HOME/.claude.json"          "archivo .claude.json"
upload_file_secret "CLAUDE_JSON_BACKUP"   "$HOME/.claude.json.backup"   "backup .claude.json"

echo "🎉 Todo listo! Los secrets están preparados."
echo "Recuerda hacer fly deploy después de subirlos."
