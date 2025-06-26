#!/bin/bash

# Test script for debugging execution order
echo "Testing execution order with delay..."

ENVOY_URL="http://localhost:10000"

echo ""
echo "=== Testing echo1 method ==="
echo "Watch the logs for execution order:"
echo "1. [EXT-PROC] Request headers (with 500ms delay)"
echo "2. [EXT-PROC] Request body"
echo "3. [WASM] Request processing" 
echo "4. Echo1 server response"
echo ""

curl -X POST $ENVOY_URL \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "echo1",
    "params": {
      "message": "Testing execution order!"
    }
  }' \
  -v

echo ""
echo ""
echo "=== Testing echo2 method ==="

curl -X POST $ENVOY_URL \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "echo2",
    "params": {
      "message": "Testing execution order 2!"
    }
  }' \
  -v

echo ""
echo "Test complete! Check the logs to see execution order." 