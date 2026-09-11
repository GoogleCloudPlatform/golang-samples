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

// [START storage_control_create_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// createAnywhereCache creates an anywhere cache for the given bucket and zone.
func createAnywhereCache(w io.Writer, bucketName, zone string) error {
	// bucketName := "bucket-name"
	// zone := "us-central1-a"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*20)
	defer cancel()

	req := &controlpb.CreateAnywhereCacheRequest{
		Parent: fmt.Sprintf("projects/_/buckets/%v", bucketName),
		AnywhereCache: &controlpb.AnywhereCache{
			Zone: zone,
		},
	}

	// Start a create/update operation and block until it completes.
	// Real applications may want to setup a callback, wait on a coroutine, or poll until it completes.
	op, err := client.CreateAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateAnywhereCache(%q): %w", zone, err)
	}

	anywhereCache, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Created anywhere cache: %v\n", anywhereCache.GetName())
	return nil
}

// [END storage_control_create_anywhere_cache]
