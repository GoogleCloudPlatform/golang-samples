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

	var conn *grpc.ClientConn
	testutil.Retry(t, 10, 20*time.Second, func(r *testutil.R) {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs: certPool,
		})))
		if err != nil {
			// Calls t.Fail() after the last attempt, code dependent on
			// a successful connection is safe because it will not be called.
			r.Errorf("grpc.Dial %s: %v", addr, err)
		}
	})
	defer conn.Close()

	var n uint32 = 3
	req := &pb.Request{DurationSecs: n}
	resp, err := pb.NewTimeServiceClient(conn).StreamTime(context.Background(), req)
	if err != nil {
		t.Fatalf("rpc StreamTime: %v", err)
	}

	var recvMsgs int
	var recvFailures int
	for {
		_, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			recvFailures++
			if recvFailures < 5 {
				t.Logf("rpc StreamTime.Recv: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				t.Fatalf("rpc StreamTime.Recv: %v", err)
			}
		}
		recvMsgs++
	}

	if recvMsgs != int(n) {
		t.Errorf("received %d messages, expected %d", recvMsgs, req.DurationSecs)
	}
}
