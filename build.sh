#!/bin/bash

echo "Building gRPC microservice..."

# Build server
echo "Building server..."
go build -o bin/grpc-server cmd/server/main.go
if [ $? -eq 0 ]; then
    echo "✓ Server built successfully"
else
    echo "✗ Server build failed"
    exit 1
fi

# Build client
echo "Building client..."
go build -o bin/grpc-client cmd/client/main.go
if [ $? -eq 0 ]; then
    echo "✓ Client built successfully"
else
    echo "✗ Client build failed"
    exit 1
fi

echo "All builds completed successfully!"
echo "Binaries created in bin/ directory:"
ls -la bin/