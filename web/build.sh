#!/bin/bash
GOOS=js GOARCH=wasm go build -o assets/main.wasm ./wasm/...