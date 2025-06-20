#!/bin/bash

set -e

echo "Running code quality checks..."

echo "Formatting code..."
gofmt -s -w .
goimports -w .

echo "Running linter..."
golangci-lint run --timeout=5m

echo "Running tests..."
go test -v -race -coverprofile=coverage.out ./...

echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "All checks passed!"
echo "Coverage report: coverage.html"
