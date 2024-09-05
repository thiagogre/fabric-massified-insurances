#!/bin/bash

set -e

echo "Running tests..."
go test -v ./... || {
    echo "Tests failed"
    exit 1
}
