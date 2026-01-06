#!/bin/bash

# ==========================================================
# CONFIGURATION
# ==========================================================
REPO_OWNER="Sarthak160"  
REPO_NAME="filepecker"    
BIN_NAME="filepecker"
INSTALL_DIR="/usr/local/bin"
# ==========================================================

set -e # Exit immediately if a command exits with a non-zero status

# 1. Detect OS & Arch
OS="$(uname -s)"
ARCH="$(uname -m)"
EXT=""

case "${OS}" in
    Linux*)     OS='linux';;
    Darwin*)    OS='darwin';;
    CYGWIN*|MINGW*|MSYS*) OS='windows'; EXT='.exe';;
    *)          echo "Error: Unsupported OS: ${OS}"; exit 1;;
esac

case "${ARCH}" in
    x86_64)    ARCH='amd64';;
    arm64|aarch64) ARCH='arm64';;
    *)         echo "Error: Unsupported Architecture: ${ARCH}"; exit 1;;
esac

TARGET_BINARY="${BIN_NAME}-${OS}-${ARCH}${EXT}"
DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/${TARGET_BINARY}"

# 2. Download
echo "⬇️  Downloading ${BIN_NAME} for ${OS}/${ARCH}..."
tmp_dir=$(mktemp -d)
curl -fsSL "$DOWNLOAD_URL" -o "${tmp_dir}/${BIN_NAME}"

# 3. Make Executable
chmod +x "${tmp_dir}/${BIN_NAME}"

# 4. Install
echo "📦 Installing to ${INSTALL_DIR}..."

# Check if we have write access to INSTALL_DIR, otherwise use sudo
if [ -w "$INSTALL_DIR" ]; then
    mv "${tmp_dir}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
else
    sudo mv "${tmp_dir}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
fi

# 5. Cleanup
rm -rf "$tmp_dir"

echo "✅ Success! Run '${BIN_NAME}' to get started."