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

// [START storage_control_create_rapid_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// createRapidCache creates a Rapid Cache for the bucket in the given zone.
func createRapidCache(w io.Writer, bucket, zone string) error {
	// bucket := "bucket-name"
	// The zone must be in the same region as the bucket.
	// zone := "us-central1-a"

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
	parent := fmt.Sprintf("projects/_/buckets/%v", bucket)
	// The Rapid Cache ID must be the same as the zone.
	name := fmt.Sprintf("%v/rapidCaches/%v", parent, zone)

	req := &controlpb.CreateRapidCacheRequest{
		Parent: parent,
		RapidCache: &controlpb.RapidCache{
			Name: name,
			Zone: zone,
			// Rapid Cache Ultra requires a bucket with hierarchical namespace
			// enabled.
			CacheType: "rapid-cache-ultra",
			// How long a cache entry can live. The TTL must be in whole
			// seconds, between 1 hour and 7 days.
			Ttl: durationpb.New(24 * time.Hour),
		},
	}
	// Creating a Rapid Cache is a long-running operation.
	op, err := client.CreateRapidCache(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateRapidCache(%q): %w", parent, err)
	}

	// Wait for the long-running operation to complete.
	rc, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for CreateRapidCache: %w", err)
	}

	fmt.Fprintf(w, "Created Rapid Cache: %v\n", rc.GetName())
	fmt.Fprintf(w, "State: %v\n", rc.GetState())
	return nil
}

// [END storage_control_create_rapid_cache]
