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

// [START storage_control_get_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// getAnywhereCache gets an Anywhere Cache instance.
func getAnywhereCache(w io.Writer, bucket, anywhereCacheID string) error {
	// bucket := "bucket-name"
	// anywhereCacheID := "anywhere-cache-id"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req := &controlpb.GetAnywhereCacheRequest{
		Name: fmt.Sprintf("projects/_/buckets/%v/anywhereCaches/%v", bucket, anywhereCacheID),
	}
	anywhereCache, err := client.GetAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("GetAnywhereCache: %w", err)
	}

	fmt.Fprintf(w, "Anywhere Cache: %v\n", anywhereCache.GetName())
	fmt.Fprintf(w, "State: %v\n", anywhereCache.GetState())
	fmt.Fprintf(w, "Admission Policy: %v\n", anywhereCache.GetAdmissionPolicy())
	return nil
}

// [END storage_control_get_anywhere_cache]
