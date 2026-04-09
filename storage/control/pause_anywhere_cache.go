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

// [START storage_control_pause_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// pauseAnywhereCache pauses an anywhere cache.
func pauseAnywhereCache(w io.Writer, bucket, anywhereCacheID string) error {
	// bucket := "bucket-name"
	// anywhereCacheID := "us-central1-a"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("storagecontrol.NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	name := fmt.Sprintf("projects/_/buckets/%s/anywhereCaches/%s", bucket, anywhereCacheID)
	req := &controlpb.PauseAnywhereCacheRequest{
		Name: name,
	}

	anywhereCache, err := client.PauseAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("PauseAnywhereCache: %w", err)
	}

	fmt.Fprintf(w, "Paused anywhere cache: %v\n", anywhereCache.GetName())
	return nil
}

// [END storage_control_pause_anywhere_cache]
