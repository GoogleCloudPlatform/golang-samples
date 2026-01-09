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

// Package client provides some basic snippet examples for working with API clients
// using the preview BigQuery Cloud Client Library.
package client

import (
	"context"
	"io"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestClients(t *testing.T) {
	testutil.SystemTest(t)
	ctx := context.Background()

	names := []string{"gRPC", "REST"}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := basicClientUsage(ctx, io.Discard, name == "gRPC")
			if err != nil {
				t.Errorf("basicClientUsage: %v", err)
			}
		})
	}
}
