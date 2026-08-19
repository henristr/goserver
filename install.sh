#!/bin/sh
set -e

REPO="henristr/goserver"
BIN="goserver"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Nicht unterstützte Architektur: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    linux)
        OS="Linux"
        ;;
    darwin)
        OS="Darwin"
        ;;
    *)
        echo "Nicht unterstütztes Betriebssystem: $OS"
        echo "Für Windows bitte PowerShell verwenden."
        exit 1
        ;;
esac

RELEASE=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")

URL=$(echo "$RELEASE" \
    | grep -o 'https://[^"]*' \
    | grep "${OS}_${ARCH}\.tar\.gz" \
    | head -n 1)

if [ -z "$URL" ]; then
    echo "Kein passendes Release für ${OS}_${ARCH} gefunden."
    exit 1
fi

echo "Lade $URL ..."

TMP="/tmp/$BIN.tar.gz"

curl -fsSL "$URL" -o "$TMP"

echo "Installiere $BIN..."

sudo tar -xzf "$TMP" -C /usr/local/bin "$BIN"

rm "$TMP"

echo ""
echo "$BIN erfolgreich installiert!"
echo "Installiert nach: /usr/local/bin/$BIN"