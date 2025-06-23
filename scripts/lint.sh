#!/bin/sh

# This script runs the linter on the codebase.
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run