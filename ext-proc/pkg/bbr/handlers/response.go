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
	"log"

	eppb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

// HandleResponseHeaders handles response headers.
func (s *Server) HandleResponseHeaders(headers *eppb.HttpHeaders) ([]*eppb.ProcessingResponse, error) {
	log.Println("[EXT-PROC] Processing response headers...")

	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_ResponseHeaders{
				ResponseHeaders: &eppb.HeadersResponse{},
			},
		},
	}, nil
}

// HandleResponseBody handles response bodies.
func (s *Server) HandleResponseBody(body *eppb.HttpBody) ([]*eppb.ProcessingResponse, error) {
	log.Printf("[EXT-PROC] Processing response body... (size: %d, end_of_stream: %t)",
		len(body.GetBody()), body.GetEndOfStream())

	// Log the response body content if it's not too large
	if len(body.GetBody()) > 0 && len(body.GetBody()) < 1000 {
		log.Printf("[EXT-PROC] Response body content: %s", string(body.GetBody()))
	}

	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_ResponseBody{
				ResponseBody: &eppb.BodyResponse{},
			},
		},
	}, nil
}

// HandleResponseTrailers handles response trailers.
func (s *Server) HandleResponseTrailers(trailers *eppb.HttpTrailers) ([]*eppb.ProcessingResponse, error) {
	log.Println("[EXT-PROC] Processing response trailers...")

	return []*eppb.ProcessingResponse{
		{
			Response: &eppb.ProcessingResponse_ResponseTrailers{
				ResponseTrailers: &eppb.TrailersResponse{},
			},
		},
	}, nil
}
