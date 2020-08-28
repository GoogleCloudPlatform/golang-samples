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

package cloudruntests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming/pkg/api/v1"
)

// TestGRPCServerStreamingService is an end-to-end test that confirms the image builds, deploys and runs on
// Cloud Run and can stream messages from server.
func TestGRPCServerStreamingService(t *testing.T) {
	t.Skip("Feature not yet available (https://github.com/GoogleCloudPlatform/golang-samples/issues/1628)")
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("grpc-server-streaming", tc.ProjectID)
	service.Dir = "../grpc-server-streaming"
	service.AllowUnauthenticated = true

	if err := service.Build(); err != nil {
		t.Fatalf("Service.Build: %v", err)
	}
	if err := service.Deploy(); err != nil {
		t.Fatalf("Service.Deploy: %v", err)
	}
	defer service.Clean()

	svcURL, err := service.ParsedURL()
	if err != nil {
		t.Fatalf("Service.ParsedURL: %v", err)
	}
	addr := svcURL.Host + ":443"
	certPool, err := x509.SystemCertPool()
	if err != nil {
		t.Fatalf("x509.SystemCertPool: %v", err)
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})))
	if err != nil {
		t.Fatalf("grpc.Dial %s: %v", addr, err)
	}
	defer conn.Close()

	var n uint32 = 3
	req := &pb.Request{DurationSecs: n}
	resp, err := pb.NewTimeServiceClient(conn).StreamTime(context.Background(), req)
	if err != nil {
		t.Fatalf("rpc StreamTime: %v", err)
	}

	recvSeconds := make(map[time.Time]bool)
	recvMsgs := 0
	for {
		_, err := resp.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("rpc StreamTime.Recv: %v", err)
		}

		recvMsgs++
		recvAt := time.Now().Truncate(time.Second)
		recvSeconds[recvAt] = true
	}

	if recvMsgs != int(n) {
		t.Errorf("received %d messages, expected %d", recvMsgs, req.DurationSecs)
	}
	if len(recvSeconds) != int(n) {
		var keys []string
		for k := range recvSeconds {
			keys = append(keys, k.Format(time.RFC3339))
		}
		t.Errorf("received messages at %d distinct seconds [%s], expected %d different seconds", len(recvSeconds), strings.Join(keys, ", "), req.DurationSecs)
	}
}
