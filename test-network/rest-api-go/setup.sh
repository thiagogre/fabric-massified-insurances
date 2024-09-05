#!/bin/bash

set -e

echo "Setting up Go project..."

echo "Getting Go modules..."
go mod tidy

./run_tests.sh

echo "Go project setup complete!"
