## Body Base Routing (BBR) from Gateway API Inference Extension

## Overview

This repository serves as an experimentation environment for testing **[body-based routing](https://github.com/kubernetes-sigs/gateway-api-inference-extension/tree/main/pkg/bbr)** using Envoy proxy with custom filters. The primary goals are:

1. **Test JSON-RPC body parsing and routing** - Route requests based on the `method` field in JSON-RPC request bodies
2. **Validate execution order** - Understand the sequence of filter processing phases (headers, body, etc.)
3. **Debug filter interactions** - Observe how ext-proc and WASM filters work together



### Run Env

```sh
make up
```

#### Run Traffic

```sh
./test-execution-order.sh
```

Check the output of the script, and the logs from the running containers.
You should see something like this:

```
ext-proc  | 2025/06/26 16:26:54 Processing new request
ext-proc  | 2025/06/26 16:26:54 Processing request with ID: 919af242-49f9-4a5f-bbd2-17d4e39d5f27
ext-proc  | 2025/06/26 16:26:54 [EXT-PROC] Starting request header processing...
ext-proc  | 2025/06/26 16:26:54 [EXT-PROC] Adding 500ms delay for debugging...
ext-proc  | 2025/06/26 16:26:54 [EXT-PROC] Completed request header processing
ext-proc  | 2025/06/26 16:26:54 Response generated: request_headers:{}
ext-proc  | 2025/06/26 16:26:54 Incoming body chunk: {
ext-proc  |     "jsonrpc": "2.0",
ext-proc  |     "id": 1,
ext-proc  |     "method": "echo1",
ext-proc  |     "params": {
ext-proc  |       "message": "Testing execution order!"
ext-proc  |     }
ext-proc  |   } (EoS: true)
ext-proc  | 2025/06/26 16:26:54 [EXT-PROC] Starting request body processing...
ext-proc  | 2025/06/26 16:26:54 [EXT-PROC] Adding 500ms delay for request body processing...
ext-proc  | 2025/06/26 16:26:55 Found JSON-RPC method: echo1
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Extracted JSON-RPC method: echo1
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Completed request body processing (non-streaming)
ext-proc  | 2025/06/26 16:26:55 Response generated: request_body:{response:{header_mutation:{set_headers:{header:{key:"x-rpc-method" raw_value:"echo1"}}} clear_route_cache:true}}
envoy     | [2025-06-26 16:26:55.314][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Starting request header processing... (headers: 11, end_of_stream: false)
envoy     | [2025-06-26 16:26:55.314][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Completed request header processing
envoy     | [2025-06-26 16:26:55.314][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Starting request body processing... (size: 129, end_of_stream: true)
envoy     | [2025-06-26 16:26:55.314][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Completed request body processing
envoy     | [2025-06-26 16:26:55.315][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Starting response header processing... (headers: 7, end_of_stream: false)
envoy     | [2025-06-26 16:26:55.315][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Completed response header processing
echo1     | 2025/06/26 16:26:55 localhost:10000 192.168.107.5:59620 "POST / HTTP/1.1" 200 26 "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)" 27.168µs
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Starting response header processing...
envoy     | [2025-06-26 16:26:55.315][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Starting response body processing... (size: 26, end_of_stream: false)
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Adding 500ms delay for response header processing...
envoy     | [2025-06-26 16:26:55.315][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Starting response body processing... (size: 0, end_of_stream: true)
envoy     | [2025-06-26 16:26:55.315][27][info][wasm] [source/extensions/common/wasm/context.cc:1195] wasm log debug_wasm_filter debug_wasm_filter_root my_vm: [WASM] Completed response body processing
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Completed response header processing
ext-proc  | 2025/06/26 16:26:55 Response generated: response_headers:{}
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Starting response body processing... (size: 26, end_of_stream: true)
ext-proc  | 2025/06/26 16:26:55 [EXT-PROC] Adding 500ms delay for response body processing...
ext-proc  | 2025/06/26 16:26:56 [EXT-PROC] Response body content: Hello from Echo Server 1!
ext-proc  | 2025/06/26 16:26:56 [EXT-PROC] Completed response body processing
ext-proc  | 2025/06/26 16:26:56 Response generated: response_body:{}
ext-proc  | 2025/06/26 16:26:56 Cannot receive stream request: rpc error: code = Canceled desc = context canceled
```

#### Clean up

```sh
make down
```
