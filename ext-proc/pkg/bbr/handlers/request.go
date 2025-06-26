/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	basepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	eppb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	extProcPb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/bbr/metrics"
)

const methodHeader = "x-rpc-method"

// JSONRPCRequest represents a simple JSON-RPC request structure
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// HandleRequestBody handles request bodies for JSON-RPC messages.
func (s *Server) HandleRequestBody(ctx context.Context, data map[string]any) ([]*eppb.ProcessingResponse, error) {
	log.Println("[EXT-PROC] Starting request body processing...")

	// Add 500ms delay for debugging execution order
	log.Println("[EXT-PROC] Adding 500ms delay for request body processing...")
	time.Sleep(500 * time.Millisecond)

	var ret []*eppb.ProcessingResponse

	requestBodyBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Try to parse as JSON-RPC request
	methodName := extractJSONRPCMethod(data)

	if methodName == "" {
		log.Println("Request body does not contain JSON-RPC method or method could not be extracted")
		if s.streaming {
			ret = append(ret, &eppb.ProcessingResponse{
				Response: &eppb.ProcessingResponse_RequestHeaders{
					RequestHeaders: &eppb.HeadersResponse{},
				},
			})
			ret = addStreamedBodyResponse(ret, requestBodyBytes)
			return ret, nil
		} else {
			ret = append(ret, &eppb.ProcessingResponse{
				Response: &eppb.ProcessingResponse_RequestBody{
					RequestBody: &eppb.BodyResponse{},
				},
			})
		}
		return ret, nil
	}

	log.Printf("[EXT-PROC] Extracted JSON-RPC method: %s", methodName)
	metrics.RecordSuccessCounter()

	if s.streaming {
		ret = append(ret, &eppb.ProcessingResponse{
			Response: &eppb.ProcessingResponse_RequestHeaders{
				RequestHeaders: &eppb.HeadersResponse{
					Response: &eppb.CommonResponse{
						ClearRouteCache: true,
						HeaderMutation: &eppb.HeaderMutation{
							SetHeaders: []*basepb.HeaderValueOption{
								{
									Header: &basepb.HeaderValue{
										Key:      methodHeader,
										RawValue: []byte(methodName),
									},
								},
							},
						},
					},
				},
			},
		})
		ret = addStreamedBodyResponse(ret, requestBodyBytes)
		log.Println("[EXT-PROC] Completed request body processing (streaming)")
		return ret, nil
	}

	log.Println("[EXT-PROC] Completed request body processing (non-streaming)")
	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_RequestBody{
				RequestBody: &eppb.BodyResponse{
					Response: &eppb.CommonResponse{
						// Necessary so that the new headers are used in the routing decision.
						ClearRouteCache: true,
						HeaderMutation: &eppb.HeaderMutation{
							SetHeaders: []*basepb.HeaderValueOption{
								{
									Header: &basepb.HeaderValue{
										Key:      methodHeader,
										RawValue: []byte(methodName),
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// extractJSONRPCMethod safely extracts the method from JSON-RPC request
func extractJSONRPCMethod(data map[string]any) string {
	// Check if this is a JSON-RPC request
	jsonrpcVal, ok := data["jsonrpc"]
	if !ok {
		log.Println("Request is not JSON-RPC format (missing jsonrpc field)")
		return ""
	}

	jsonrpcStr, ok := jsonrpcVal.(string)
	if !ok || jsonrpcStr != "2.0" {
		log.Println("Request is not JSON-RPC 2.0 format")
		return ""
	}

	// Extract method field
	methodVal, ok := data["method"]
	if !ok {
		log.Println("JSON-RPC request missing method field")
		return ""
	}

	methodStr, ok := methodVal.(string)
	if !ok {
		log.Println("JSON-RPC method is not a string")
		return ""
	}

	log.Printf("Found JSON-RPC method: %s", methodStr)
	return methodStr
}

func addStreamedBodyResponse(responses []*eppb.ProcessingResponse, requestBodyBytes []byte) []*eppb.ProcessingResponse {
	return append(responses, &extProcPb.ProcessingResponse{
		Response: &extProcPb.ProcessingResponse_RequestBody{
			RequestBody: &extProcPb.BodyResponse{
				Response: &extProcPb.CommonResponse{
					BodyMutation: &extProcPb.BodyMutation{
						Mutation: &extProcPb.BodyMutation_StreamedResponse{
							StreamedResponse: &extProcPb.StreamedBodyResponse{
								Body:        requestBodyBytes,
								EndOfStream: true,
							},
						},
					},
				},
			},
		},
	})
}

// HandleRequestHeaders handles request headers.
func (s *Server) HandleRequestHeaders(headers *eppb.HttpHeaders) ([]*eppb.ProcessingResponse, error) {
	log.Println("[EXT-PROC] Starting request header processing...")

	// Add 500ms delay for debugging execution order
	log.Println("[EXT-PROC] Adding 500ms delay for debugging...")
	time.Sleep(500 * time.Millisecond)

	log.Println("[EXT-PROC] Completed request header processing")

	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_RequestHeaders{
				RequestHeaders: &eppb.HeadersResponse{},
			},
		},
	}, nil
}

// HandleRequestTrailers handles request trailers.
func (s *Server) HandleRequestTrailers(trailers *eppb.HttpTrailers) ([]*eppb.ProcessingResponse, error) {
	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_RequestTrailers{
				RequestTrailers: &eppb.TrailersResponse{},
			},
		},
	}, nil
}
