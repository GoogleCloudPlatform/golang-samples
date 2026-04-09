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

// [START storage_control_update_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateAnywhereCache updates an anywhere cache admission policy.
func updateAnywhereCache(w io.Writer, bucket, anywhereCacheID, admissionPolicy string) error {
	// bucket := "bucket-name"
	// anywhereCacheID := "us-central1-a"
	// admissionPolicy := "admit-on-first-miss"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("storagecontrol.NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*20)
	defer cancel()

	name := fmt.Sprintf("projects/_/buckets/%s/anywhereCaches/%s", bucket, anywhereCacheID)
	req := &controlpb.UpdateAnywhereCacheRequest{
		AnywhereCache: &controlpb.AnywhereCache{
			Name:            name,
			AdmissionPolicy: admissionPolicy,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"admission_policy"},
		},
	}

	// UpdateAnywhereCache is a long-running operation.
	// Real applications may want to setup a callback or wait for the operation to complete.
	// Blocking with .Wait(ctx) is for simplicity and real applications may prefer callbacks, coroutines, or polling.
	op, err := client.UpdateAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateAnywhereCache: %w", err)
	}

	anywhereCache, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Updated anywhere cache: %v\n", anywhereCache.GetName())
	return nil
}

// [END storage_control_update_anywhere_cache]
