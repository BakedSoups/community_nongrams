#!/bin/sh
set -eu

scripts/write-web-config.sh
go run ./cmd/genlevels

wasm_exec="$(go env GOROOT)/misc/wasm/wasm_exec.js"
if [ -f "$wasm_exec" ]; then
  cp "$wasm_exec" static/wasm_exec.js
fi

GOOS=js GOARCH=wasm go build -buildvcs=false -o static/game.wasm ./cmd/game
