#!/bin/sh
set -e

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
[ "$ARCH" = "x86_64" ] && ARCH="amd64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"

echo "Fetching latest toofan..."
URL=$(curl -s "https://api.github.com/repos/vyrx-dev/toofan/releases/latest" | grep "browser_download_url.*toofan-${OS}-${ARCH}\"" | cut -d '"' -f 4)

if [ -z "$URL" ]; then
    echo "Error: Release not found for ${OS}-${ARCH}."
    exit 1
fi

echo "Downloading from $URL..."
curl -sL "$URL" -o /tmp/toofan
chmod +x /tmp/toofan

if [ -w "/usr/local/bin" ]; then
    mv /tmp/toofan /usr/local/bin/toofan
else
    sudo mv /tmp/toofan /usr/local/bin/toofan
fi

echo "Installed to /usr/local/bin/toofan"
sleep 0.5
/usr/local/bin/toofan
