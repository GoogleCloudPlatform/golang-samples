// Copyright 2026 Google LLC
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

package control

// [START storage_control_disable_rapid_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// disableRapidCache starts disabling a Rapid Cache. It does not wait for the
// long-running operation to complete.
func disableRapidCache(w io.Writer, bucket, cacheID string) error {
	// bucket := "bucket-name"
	// The Rapid Cache ID is the zone of the cache.
	// cacheID := "us-central1-a"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// The storage bucket path uses the global access pattern, in which the "_"
	// denotes this bucket exists in the global namespace.
	name := fmt.Sprintf("projects/_/buckets/%v/rapidCaches/%v", bucket, cacheID)

	req := &controlpb.DisableRapidCacheRequest{
		Name: name,
	}
	// Disabling a Rapid Cache starts a long-running operation. The cache state
	// changes to "disabled" within seconds, but the operation only completes
	// once the cache has been fully removed, which can take more than an hour.
	// This sample does not wait for the operation to complete; use
	// GetRapidCache to check the cache state.
	op, err := client.DisableRapidCache(ctx, req)
	if err != nil {
		return fmt.Errorf("DisableRapidCache(%q): %w", name, err)
	}

	fmt.Fprintf(w, "Disabling Rapid Cache: %v\n", name)
	fmt.Fprintf(w, "Operation: %v\n", op.Name())
	return nil
}

// [END storage_control_disable_rapid_cache]
