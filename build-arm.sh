#!/bin/bash

BINFILE=watchtower
if [ -n "$MSYSTEM" ]; then
    BINFILE=watchtower.exe
fi
echo "Before git describe"
git describe --tags
echo "After git describe"
VERSION=$(git describe --tags)
echo "Building $VERSION..."
GOOS=linux GOARCH=arm GOARM=7 go build -o $BINFILE -ldflags "-X github.com/meibye/watchtower/internal/meta.Version=$VERSION"
