#!/bin/bash

# directory to store binaries
mkdir -p build/osx-intel
mkdir -p build/osx-m1

# build for osx-intel
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o build/osx-intel/srgps main.go

# build for osx-m1
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o build/osx-m1/srgps main.go

echo "Build complete."
