#!/bin/bash

ARGS=$1

set -e

echo "Running unit tests..."
go test $ARGS $(go list ./... | grep -v '/tests/integration') || {
    echo "Tests failed"
    exit 1
}
