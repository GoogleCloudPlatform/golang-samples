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

// [START storage_control_update_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateAnywhereCache updates an Anywhere Cache.
func updateAnywhereCache(w io.Writer, bucket, zone, admissionPolicy string) error {
	// bucket := "bucket-name"
	// zone := "us-central1-a"
	// admissionPolicy := "admit-on-second-miss"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*20)
	defer cancel()

	cacheName := fmt.Sprintf("projects/_/buckets/%v/anywhereCaches/%v", bucket, zone)
	req := &controlpb.UpdateAnywhereCacheRequest{
		AnywhereCache: &controlpb.AnywhereCache{
			Name:            cacheName,
			AdmissionPolicy: admissionPolicy,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"admission_policy"},
		},
	}
	op, err := client.UpdateAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateAnywhereCache: %w", err)
	}

	cache, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Updated Anywhere Cache: %v\n", cache.Name)
	return nil
}

// [END storage_control_update_anywhere_cache]
