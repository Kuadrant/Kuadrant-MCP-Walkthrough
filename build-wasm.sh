#!/bin/bash

# Build the WASM filter
echo "Building WASM filter..."
cargo build --target wasm32-wasip1 --release

if [ $? -eq 0 ]; then
    echo "WASM filter built successfully!"
    echo "Output: target/wasm32-wasip1/release/debug_wasm_filter.wasm"
else
    echo "Failed to build WASM filter"
    exit 1
fi 