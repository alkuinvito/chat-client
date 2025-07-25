#!/bin/bash

cd "$(dirname "$0")/../../.." || exit 1

export CC=x86_64-w64-mingw32-gcc
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1

fyne package --os windows --icon assets/Icon.png --exe tmp/chat-client.exe
