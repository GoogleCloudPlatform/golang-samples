// Copyright 2021 Google LLC
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
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/http2"
)

// TestGRPCServerStreamingService is an end-to-end test that confirms the image builds, deploys and runs on
// Cloud Run and can stream messages from server.
func TestHTTP2Server(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("h2c", tc.ProjectID)
	service.Dir = "../h2c"
	service.AllowUnauthenticated = true
	service.AsBuildpack = true
	service.HTTP2 = true

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

	h2Client := &http.Client{Transport: &http2.Transport{}}
	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		resp, err := h2Client.Get(svcURL.String())

		if err != nil {
			r.Errorf("http2.Get failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			r.Errorf("http2.Get: unexpected response status: %s", resp.Status)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			r.Errorf("resp.Body.Read failed: %v", err)
		}
		if expected, got := "This request is served over HTTP/2.0 protocol.", string(b); !strings.Contains(got, expected) {
			r.Errorf("response body doesn't contain %q; got=%q", expected, got)
		}
	})
}
