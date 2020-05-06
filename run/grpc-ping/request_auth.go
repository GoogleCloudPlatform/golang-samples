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

// [START run_grpc_request_auth]

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/compute/metadata"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
)

// pingRequestWithAuth mints a new ID Token with the compute metadata server for each request.
// This token has a 1 hour expiry and should be reused.
// This function will only work on Google Cloud with an available compute metadata server and compute identity.
// audience is the auto-assigned URL of the Cloud Run service (do not include port number)
func pingRequestWithAuth(conn *grpc.ClientConn, p *pb.Request, audience string) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an ID Token as shown in service-to-service authentication in the documentation.
	// https://cloud.google.com/run/docs/authenticating/service-to-service
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", audience)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("metadata.Get: failed to query id_token: %v", err)
	}
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+idToken)

	client := pb.NewPingServiceClient(conn)
	return client.Send(ctx, p)
}

// [END run_grpc_request_auth]
