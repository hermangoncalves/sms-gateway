#!/usr/bin/env bash
set -e

REPO="hermangoncalves/sms-gateway"
VERSION="${1:-latest}"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="sms-gateway"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    armv7l) ARCH="armv7" ;;
    *) echo "❌ Arquitetura não suportada: $ARCH"; exit 1 ;;
esac

echo "📦 Instalando $BINARY_NAME ($OS-$ARCH), versão: $VERSION"

# Resolve versão mais recente se for "latest"
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -Po '"tag_name": "\\K(.*?)(?=")')
fi

if [ -z "$VERSION" ]; then
    echo "❌ Não foi possível obter a versão"
    exit 1
fi

echo "➡️ Versão: $VERSION"

# Monta URL de download
ASSET="$BINARY_NAME-$OS-$ARCH"
URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET"

TMP_FILE=$(mktemp)
echo "⬇️  Baixando: $URL"
curl -sL "$URL" -o "$TMP_FILE"

# Instala
if [ ! -w "$INSTALL_DIR" ]; then
    echo "⚠️ Sem permissão para escrever em $INSTALL_DIR, instalando em \$HOME/.local/bin"
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "✅ Instalado em $INSTALL_DIR/$BINARY_NAME"

# Teste
if command -v $BINARY_NAME >/dev/null 2>&1; then
    echo "🚀 Execução bem-sucedida: $($BINARY_NAME --help || true)"
else
    echo "⚠️ Adicione $INSTALL_DIR ao seu PATH"
fi
