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

// [START storage_control_list_anywhere_caches]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/api/iterator"
)

// listAnywhereCaches lists anywhere caches for a bucket.
func listAnywhereCaches(w io.Writer, bucket string) error {
	// bucket := "bucket-name"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("storagecontrol.NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req := &controlpb.ListAnywhereCachesRequest{
		Parent: fmt.Sprintf("projects/_/buckets/%s", bucket),
	}

	it := client.ListAnywhereCaches(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListAnywhereCaches: %w", err)
		}
		fmt.Fprintf(w, "Anywhere cache: %v\n", resp.GetName())
	}

	return nil
}

// [END storage_control_list_anywhere_caches]
