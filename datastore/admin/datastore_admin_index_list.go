// Copyright 2020 Google LLC
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

package samples

// [START datastore_admin_index_list]
import (
	"context"
	"fmt"
	"io"

	admin "cloud.google.com/go/datastore/admin/apiv1"
	"cloud.google.com/go/datastore/admin/apiv1/adminpb"
	"google.golang.org/api/iterator"
)

// indexList lists the indexes.
func indexList(w io.Writer, projectID string) ([]*adminpb.Index, error) {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := admin.NewDatastoreAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("admin.NewDatastoreAdminClient: %w", err)
	}
	defer client.Close()

	req := &adminpb.ListIndexesRequest{
		ProjectId: projectID,
	}
	it := client.ListIndexes(ctx, req)
	var indices []*adminpb.Index
	for {
		index, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ListIndexes: %w", err)
		}
		indices = append(indices, index)
		fmt.Fprintf(w, "Got index: %v\n", index.IndexId)
	}

	fmt.Fprintf(w, "Got lists of indexes\n")
	return indices, nil
}

// [END datastore_admin_index_list]
