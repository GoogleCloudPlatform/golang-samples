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

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/api/option"
)

// pauseAnywhereCache pauses an Anywhere Cache.
func pauseAnywhereCache(w io.Writer, bucket, zone string) error {
	// bucket := "bucket-name"
	// zone := "us-central1-f"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx, option.WithQuotaProject(""))
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req := &controlpb.PauseAnywhereCacheRequest{
		Name: fmt.Sprintf("projects/_/buckets/%v/anywhereCaches/%v", bucket, zone),
	}
	resp, err := client.PauseAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("PauseAnywhereCache: %w", err)
	}

	fmt.Fprintf(w, "paused anywhere cache with path %q", resp.Name)
	return nil
}

// [END storage_control_pause_anywhere_cache]
