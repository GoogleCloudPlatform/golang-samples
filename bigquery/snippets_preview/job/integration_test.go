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

// Package job provides some basic snippet examples for working with job
// metadata using the preview BigQuery Cloud Client Library.
package job

import (
	"context"
	"io"
	"testing"
	"time"

	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const testTimeout = 30 * time.Second

func TestJobSnippet(t *testing.T) {
	tc := testutil.SystemTest(t)
	names := []string{"gRPC", "REST"}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
			defer cancel()
			// Setup client.
			var client *apiv2_client.Client
			var err error
			if name == "gRPC" {
				client, err = apiv2_client.NewClient(ctx)
			} else {
				client, err = apiv2_client.NewRESTClient(ctx)
			}
			if err != nil {
				t.Fatalf("client creation failed: %v", err)
			}
			defer client.Close()
			// Create a dataset.
			projID := tc.ProjectID
			if err := listJobs(client, io.Discard, projID); err != nil {
				t.Fatalf("listJobs(%q): %v", projID, err)
			}
		})
	}
}
