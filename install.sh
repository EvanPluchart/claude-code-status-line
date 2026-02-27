#!/bin/sh
# Claude Code Status Line - Install Script
# Usage: curl -sSL https://raw.githubusercontent.com/EvanPluchart/claude-code-status-line/main/install.sh | sh

set -e

REPO="EvanPluchart/claude-code-status-line"
BINARY="claude-code-status-line"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux)  EXT="tar.gz" ;;
  darwin) EXT="tar.gz" ;;
  *)
    echo "Unsupported OS: $OS"
    echo "For Windows, use install.ps1 or download from GitHub Releases."
    exit 1
    ;;
esac

# Get latest version
LATEST=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Error: Could not determine latest version."
  exit 1
fi

URL="https://github.com/${REPO}/releases/download/v${LATEST}/${BINARY}_${LATEST}_${OS}_${ARCH}.${EXT}"

echo "Installing ${BINARY} v${LATEST} (${OS}/${ARCH})..."

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -sSL "$URL" -o "${TMP_DIR}/archive.${EXT}"
tar -xzf "${TMP_DIR}/archive.${EXT}" -C "$TMP_DIR"

if [ -w "$INSTALL_DIR" ]; then
  cp "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  sudo cp "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"

echo "Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
echo ""
echo "Run '${BINARY} init' to configure your statusline."
