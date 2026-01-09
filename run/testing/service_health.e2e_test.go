// Copyright 2025 Google LLC
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
	"io"
	"net/http"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestServiceHealth(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("service-health", tc.ProjectID)
	service.Readiness = &cloudrunci.ReadinessProbe{
		TimeoutSeconds:   1,
		PeriodSeconds:    1,
		SuccessThreshold: 1,
		FailureThreshold: 1,
		HttpGet: &cloudrunci.HTTPGetProbe{
			Path: "/are_you_ready",
			Port: 8080,
		},
	}
	service.Dir = "../service-health"
	service.AsBuildpack = true
	service.Platform.CommandFlags()

	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer func(service *cloudrunci.Service) {
		err := service.Clean()
		if err != nil {
			t.Fatalf("service.Clean %q: %v", service.Name, err)
		}
	}(service)

	resp, err := service.Request("GET", "/are_you_ready")
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}

	if got, want := string(out), "HEALTHY"; got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer c.Close()
	bucketName := os.Getenv("GOLANG_SAMPLES_PROJECT_ID") + "-" + service.Version()
	t.Logf("Deleting bucket: %s", bucketName)

	err = testutil.DeleteBucketIfExists(ctx, c, bucketName)
	if err != nil {
		t.Fatalf("testutil.DeleteBucketIfExists: %v", err)
	}
}
