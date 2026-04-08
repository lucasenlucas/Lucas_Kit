#!/bin/bash
set -e

REPO="lucasenlucas/NetScope"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  i386|i686) ARCH="386" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "=> Fetching latest release of NetScope..."
# Get the browser_download_url matching the os and arch
LATEST_URL=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"browser_download_url":' | grep -v 'checksums.txt' | grep -i "netscope_${OS}_${ARCH}" | head -n 1 | cut -d '"' -f 4)

if [ -z "$LATEST_URL" ]; then
    echo "=> Could not find a release for ${OS}_${ARCH}. Please download manually or use 'go install'."
    exit 1
fi

echo "=> Downloading from $LATEST_URL..."
TMP_DIR=$(mktemp -d)
curl -sL "$LATEST_URL" -o "$TMP_DIR/netscope.tar.gz"

echo "=> Extracting..."
tar -xzf "$TMP_DIR/netscope.tar.gz" -C "$TMP_DIR" netscope || {
    echo "Failed to extract binary."
    rm -rf "$TMP_DIR"
    exit 1
}

echo "=> Installing to /usr/local/bin (may require sudo)..."
sudo mv "$TMP_DIR/netscope" /usr/local/bin/netscope
sudo chmod +x /usr/local/bin/netscope

rm -rf "$TMP_DIR"

echo "✅ NetScope installed successfully! Run 'netscope --help' to get started."
