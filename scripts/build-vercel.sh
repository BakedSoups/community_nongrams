#!/bin/sh
set -eu

if ! command -v go >/dev/null 2>&1; then
  go_version="$(awk '$1 == "go" { print $2; exit }' go.mod)"
  go_root="${VERCEL_CACHE_DIR:-/tmp}/go-${go_version}"

  if [ ! -x "$go_root/bin/go" ]; then
    archive="/tmp/go-${go_version}.tar.gz"
    rm -rf "$go_root"
    mkdir -p "$go_root"
    curl -fsSL "https://go.dev/dl/go${go_version}.linux-amd64.tar.gz" -o "$archive"
    tar -xzf "$archive" -C "$go_root" --strip-components=1
  fi

  PATH="$go_root/bin:$PATH"
  export PATH
fi

scripts/write-web-config.sh
go run ./cmd/genlevels

wasm_exec="$(go env GOROOT)/misc/wasm/wasm_exec.js"
if [ -f "$wasm_exec" ]; then
  cp "$wasm_exec" static/wasm_exec.js
fi

GOOS=js GOARCH=wasm go build -buildvcs=false -o static/game.wasm ./cmd/game
