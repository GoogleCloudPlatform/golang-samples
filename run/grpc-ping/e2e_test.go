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

package main_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
)

// TestGRPCPingService is an end-to-end test that confirms ping responses work as expected with authentication.
// Test Cases:
// - Authenticated request from client to ping service
// - Authenticated request from client to ping service w/ recursive upstream request
func TestGRPCPingService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	// Prepare the container image for both services.
	pingService := cloudrunci.NewService("grpc-ping", tc.ProjectID)
	if err := pingService.Build(); err != nil {
		t.Fatalf("Service.Build %q: %v", pingService.Name, err)
	}

	// Deploy the ping-upstream service.
	upstreamService := cloudrunci.NewService("grpc-ping-upstream", tc.ProjectID)
	upstreamService.Image = pingService.Image

	if err := upstreamService.Deploy(); err != nil {
		t.Fatalf("Service.Deploy %q: %v", upstreamService.Name, err)
	}
	defer upstreamService.Clean()

	// Deploy the ping service.
	upstreamHost, err := upstreamService.Host()
	if err != nil {
		t.Fatalf("Service.Host %q: %v", upstreamService.Name, err)
	}
	pingService.Env = cloudrunci.EnvVars{
		"GRPC_PING_HOST": upstreamHost,
	}
	if err := pingService.Deploy(); err != nil {
		t.Fatalf("Service.Deploy %q: %v", pingService.Name, err)
	}
	defer pingService.Clean()

	// Test a gRPC Request.
	pingURL, err := pingService.ParsedURL()
	if err != nil {
		t.Fatalf("Service.ParsedURL %q: %v", pingService.Name, err)
	}

	message := "Hello Tester"
	resp, err := grpcRequest(pingURL.Host+":443", pingURL.String(), func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := pb.NewPingServiceClient(conn)
		return c.Send(ctx, &pb.Request{
			Message: message,
		})
	})
	if err != nil {
		t.Fatalf("grpcRequest (Send): %v", err)
	}
	got := resp.(*pb.Response).Pong.GetMessage()
	if want := message; got != want {
		t.Errorf("response: got %q, want %q", got, want)
	}

	// Test a gRPC request for upstream relay.
	resp, err = grpcRequest(pingURL.Host+":443", pingURL.String(), func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := pb.NewPingServiceClient(conn)
		return c.SendUpstream(ctx, &pb.Request{
			Message: message,
		})
	})
	if err != nil {
		t.Fatalf("grpcRequest (SendUpstream): %v", err)
	}

	got = resp.(*pb.Response).Pong.GetMessage()
	if want := fmt.Sprintf("%s (relayed)", message); got != want {
		t.Errorf("response: got %q, want %q", got, want)
	}
}

// grpcRequest takes a callback to issue a gRPC request.
// The returned interface{} can be coerced into the expected *pb.Response.
func grpcRequest(host string, audience string, fn func(context.Context, *grpc.ClientConn) (interface{}, error)) (interface{}, error) {
	// Create gRPC connection
	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %v", err)
	}
	defer conn.Close()

	// Create an authenticated gRPC ping request.
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	idToken, err := cloudrunci.CreateIDToken(audience)
	if err != nil {
		return nil, fmt.Errorf("cloudrunci.CreateIDToken: %v", err)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+idToken)

	return fn(ctx, conn)
}
