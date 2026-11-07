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

// [START storage_control_list_rapid_caches]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
	"google.golang.org/api/iterator"
)

// listRapidCaches lists all Rapid Caches for the bucket.
func listRapidCaches(w io.Writer, bucket string) error {
	// bucket := "bucket-name"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// The storage bucket path uses the global access pattern, in which the "_"
	// denotes this bucket exists in the global namespace.
	parent := fmt.Sprintf("projects/_/buckets/%v", bucket)

	req := &controlpb.ListRapidCachesRequest{
		Parent: parent,
	}
	// The iterator fetches additional pages of results as needed.
	it := client.ListRapidCaches(ctx, req)
	for {
		rc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListRapidCaches(%q): %w", parent, err)
		}
		fmt.Fprintf(w, "Rapid Cache: %v\n", rc.GetName())
	}

	return nil
}

// [END storage_control_list_rapid_caches]
