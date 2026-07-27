#!/bin/bash
set -e

echo "=== DevAgent CLI Setup ==="

if [ ! -f .env ]; then
    cp .env.example .env
    echo "✔ Created .env from .env.example"
    echo "  → Edit .env with your API keys"
else
    echo "✔ .env already exists"
fi

if command -v go &> /dev/null; then
    echo "✔ Go $(go version | awk '{print $3}')"
else
    echo "✖ Go not found. Install from https://go.dev/dl/"
    exit 1
fi

echo ""
echo "Building..."
go build -o devagent .
echo "✔ Built: ./devagent"

echo ""
echo "=== Setup complete ==="
echo ""
echo "Try:"
echo "  ./devagent config              # check configuration"
echo "  ./devagent ask 'hello world'   # ask the LLM"
echo "  ./devagent review --staged     # review git diff"
