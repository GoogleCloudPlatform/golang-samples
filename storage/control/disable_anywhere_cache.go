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

// [START storage_control_disable_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// disableAnywhereCache disables an Anywhere Cache instance.
func disableAnywhereCache(w io.Writer, bucket, cacheId string) error {
	// bucket := "bucket-name"
	// cacheId := "00b0c000-0000-0000-0000-000000000000"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	name := fmt.Sprintf("projects/_/buckets/%v/anywhereCaches/%v", bucket, cacheId)
	req := &controlpb.DisableAnywhereCacheRequest{
		Name: name,
	}

	cache, err := client.DisableAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("DisableAnywhereCache: %w", err)
	}

	fmt.Fprintf(w, "Disabled anywhere cache: %v\n", cache.GetName())
	return nil
}

// [END storage_control_disable_anywhere_cache]
