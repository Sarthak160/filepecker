#!/bin/bash

# Configuration
REPO="YOUR_GITHUB_USERNAME/YOUR_REPO_NAME"
APP_NAME="filepecker"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS='linux';;
    Darwin*)    OS='darwin';;
    CYGWIN*|MINGW*|MSYS*) OS='windows';;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)    ARCH='amd64';;
    arm64|aarch64) ARCH='arm64';;
    *)         echo "Unsupported Architecture: ${ARCH}"; exit 1;;
esac

# Determine extension
EXT=""
if [ "$OS" == "windows" ]; then
    EXT=".exe"
fi

# Construct Download URL (Latest Release)
BINARY_NAME="${APP_NAME}-${OS}-${ARCH}${EXT}"
URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"

echo "Detected ${OS} ${ARCH}..."
echo "Downloading ${APP_NAME} from ${URL}..."

# Download
curl -L -o "${APP_NAME}" "${URL}"

# Make executable
chmod +x "${APP_NAME}"

# Move to install directory (requires sudo)
echo "Installing to ${INSTALL_DIR} (requires password)..."
sudo mv "${APP_NAME}" "${INSTALL_DIR}/${APP_NAME}"

echo "Success! Run '${APP_NAME}' to start."