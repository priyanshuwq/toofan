#!/bin/sh
set -e

# Prevent execution if this script was only partially downloaded
{
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    [ "$ARCH" = "x86_64" ] && ARCH="amd64"
    [ "$ARCH" = "aarch64" ] && ARCH="arm64"

    if ! command -v curl >/dev/null 2>&1; then
        echo "Error: curl is required to install toofan."
        exit 1
    fi

    echo "Fetching latest toofan..."
    URL="https://github.com/vyrx-dev/toofan/releases/latest/download/toofan-${OS}-${ARCH}"

    if ! curl -sfL "$URL" -o /tmp/toofan; then
        echo "Error: Release not found for ${OS}-${ARCH}."
        exit 1
    fi
    chmod +x /tmp/toofan

    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    mv /tmp/toofan "$INSTALL_DIR/toofan"
    echo "Installed to $INSTALL_DIR/toofan"

    # Add alias for toofan binary to user's shell config
    case "$SHELL" in
        */bash) SHELL_CONFIG="$HOME/.bashrc" ;;
        */zsh)  SHELL_CONFIG="$HOME/.zshrc" ;;
        */fish) SHELL_CONFIG="$HOME/.config/fish/config.fish"; mkdir -p "$(dirname "$SHELL_CONFIG")" ;;
        *)      SHELL_CONFIG="" ;;
    esac

    if [ -n "$SHELL_CONFIG" ]; then
        if ! grep -qF "alias toofan=" "$SHELL_CONFIG" 2>/dev/null; then
            printf '\n# toofan\nalias toofan="$HOME/.local/bin/toofan"\n' >> "$SHELL_CONFIG"
        fi
    else
        echo "Add alias toofan=\"\$HOME/.local/bin/toofan\" to your shell config."
    fi

    export PATH="$INSTALL_DIR:$PATH"
    sleep 0.5
    toofan
}
