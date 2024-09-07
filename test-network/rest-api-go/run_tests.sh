#!/bin/bash

ARGS=$1

set -e

echo "Running tests..."
go test $ARGS ./... || {
    echo "Tests failed"
    exit 1
}
