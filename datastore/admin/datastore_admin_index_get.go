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

// [START datastore_admin_index_get]
import (
	"context"
	"fmt"
	"io"

	admin "cloud.google.com/go/datastore/admin/apiv1"
	"cloud.google.com/go/datastore/admin/apiv1/adminpb"
)

// indexGet gets an index.
func indexGet(w io.Writer, projectID, indexID string) (*adminpb.Index, error) {
	// projectID := "my-project-id"
	// indexID := "my-index"
	ctx := context.Background()
	client, err := admin.NewDatastoreAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("admin.NewDatastoreAdminClient: %w", err)
	}
	defer client.Close()

	req := &adminpb.GetIndexRequest{
		ProjectId: projectID,
		IndexId:   indexID,
	}
	index, err := client.GetIndex(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client.GetIndex: %w", err)
	}

	fmt.Fprintf(w, "Got index: %v\n", index.IndexId)
	return index, nil
}

// [END datastore_admin_index_get]
