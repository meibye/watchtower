#!/bin/bash

VERSION=$(git describe --tags)
echo "Building $VERSION..."
GOOS=linux GOARCH=arm GOARM=7 go build -o watchtower -ldflags "-X github.com/meibye/watchtower/internal/meta.Version=$VERSION"
