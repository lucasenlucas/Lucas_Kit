#!/bin/bash

APP_NAME="netscope"
VERSION="9.4.2029"

echo "Building $APP_NAME v$VERSION..."

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o builds/${APP_NAME}_linux_amd64 ./cmd/netscope
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.version=$VERSION" -o builds/${APP_NAME}_linux_arm64 ./cmd/netscope
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o builds/${APP_NAME}_windows_amd64.exe ./cmd/netscope
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.version=$VERSION" -o builds/${APP_NAME}_darwin_amd64 ./cmd/netscope
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X main.version=$VERSION" -o builds/${APP_NAME}_darwin_arm64 ./cmd/netscope

echo "Build complete. Binaries are in the builds/ directory."
