#!/bin/bash

# Build the Go binary
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# Package the binary into a zip file
zip function.zip bootstrap

# Clean up
rm bootstrap