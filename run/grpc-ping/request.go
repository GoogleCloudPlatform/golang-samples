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

// [START cloudrun_grpc_request]

import (
	"context"
	"time"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
	"google.golang.org/grpc"
)

// pingRequest sends a new gRPC ping request to the server configured in the connection.
func pingRequest(conn *grpc.ClientConn, p *pb.Request) (*pb.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := pb.NewPingServiceClient(conn)
	return client.Send(ctx, p)
}

// [END cloudrun_grpc_request]

// PingRequest creates a new gRPC request to the upstream ping gRPC service.
func PingRequest(conn *grpc.ClientConn, p *pb.Request, url string, authenticated bool) (*pb.Response, error) {
	if authenticated {
		return pingRequestWithAuth(conn, p, url)
	}
	return pingRequest(conn, p)
}
