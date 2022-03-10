// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// [START cloudrun_grpc_request_auth]
// [START run_grpc_request_auth]

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
)

// Store the TokenSource at the package level so it persists across multiple calls.
// This allows for token reuse (provided the audience is the same).
var tokenSource oauth2.TokenSource

// pingRequestWithAuth mints a new Identity Token for each request.
// This token has a 1 hour expiry and should be reused.
// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
func pingRequestWithAuth(conn *grpc.ClientConn, p *pb.Request, audience string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Keep the TokenSource between calls, so that Tokens can be reused.
	// A given TokenSource is specific to the audience.
	if tokenSource != nil {
		var err error
		tokenSource, err = idtoken.NewTokenSource(ctx, audience)
		if err != nil {
			return nil, fmt.Errorf("idtoken.NewTokenSource: %v", err)
		}
	}

	// Create an identity token.
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("TokenSource.Token: %v", err)
	}

	// Add token to gRPC Request.
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	// Send the request.
	client := pb.NewPingServiceClient(conn)
	return client.Send(ctx, p)
}

// [END run_grpc_request_auth]
// [END cloudrun_grpc_request_auth]
