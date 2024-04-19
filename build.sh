#!/usr/bin/env bash
set -euo pipefail

OS="${GOOS:-}"
ARCH="${GOARCH:-}"

if [[ "$OS" == "" ]]; then
  OS=$(uname | awk '{print tolower($0)}')
fi
if [[ "$ARCH" == "" ]]; then
  ARCH=$(uname -m | awk '{print tolower($0)}')
fi

OUT_PATH="bin/$OS-$ARCH/webview"

cd webview
go build -o "../$OUT_PATH" .

echo "Wrote $OUT_PATH"
