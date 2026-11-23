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

// [START storage_control_update_rapid_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateRapidCache updates the time to live (TTL) of a Rapid Cache.
func updateRapidCache(w io.Writer, bucket, cacheID string, ttl time.Duration) error {
	// bucket := "bucket-name"
	// The Rapid Cache ID is the zone of the cache.
	// cacheID := "us-central1-a"
	// The TTL must be in whole seconds, between 1 hour and 7 days.
	// ttl := 48 * time.Hour

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	// The storage bucket path uses the global access pattern, in which the "_"
	// denotes this bucket exists in the global namespace.
	name := fmt.Sprintf("projects/_/buckets/%v/rapidCaches/%v", bucket, cacheID)

	req := &controlpb.UpdateRapidCacheRequest{
		RapidCache: &controlpb.RapidCache{
			Name: name,
			Ttl:  durationpb.New(ttl),
		},
		// Only the fields listed in the update mask are changed.
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"ttl"},
		},
	}
	// Updating a Rapid Cache is a long-running operation.
	op, err := client.UpdateRapidCache(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateRapidCache(%q): %w", name, err)
	}

	// Wait for the long-running operation to complete.
	rc, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for UpdateRapidCache: %w", err)
	}

	fmt.Fprintf(w, "Updated Rapid Cache: %v\n", rc.GetName())
	fmt.Fprintf(w, "TTL: %v seconds\n", rc.GetTtl().GetSeconds())
	return nil
}

// [END storage_control_update_rapid_cache]
