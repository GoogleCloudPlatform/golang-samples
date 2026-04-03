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

package control

// [START storage_control_create_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// createAnywhereCache creates an Anywhere Cache in the bucket.
func createAnywhereCache(w io.Writer, bucket, zone string) error {
	// bucket := "bucket-name"
	// zone := "us-central1-a"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*20)
	defer cancel()

	req := &controlpb.CreateAnywhereCacheRequest{
		Parent: fmt.Sprintf("projects/_/buckets/%v", bucket),
		AnywhereCache: &controlpb.AnywhereCache{
			Zone: zone,
		},
	}
	op, err := client.CreateAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateAnywhereCache: %w", err)
	}

	cache, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Created Anywhere Cache: %v\n", cache.Name)
	return nil
}

// [END storage_control_create_anywhere_cache]
