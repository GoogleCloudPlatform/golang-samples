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

// [START datastore_admin_entities_import]
import (
	"context"
	"fmt"
	"io"

	admin "cloud.google.com/go/datastore/admin/apiv1"
	"cloud.google.com/go/datastore/admin/apiv1/adminpb"
)

// entitiesImport imports entities into Datastore.
func entitiesImport(w io.Writer, projectID, inputURL string) error {
	// projectID := "project-id"
	// inputURL := "gs://bucket-name/overall-export-metadata-file"
	ctx := context.Background()
	client, err := admin.NewDatastoreAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("admin.NewDatastoreAdminClient: %w", err)
	}
	defer client.Close()

	req := &adminpb.ImportEntitiesRequest{
		ProjectId: projectID,
		InputUrl:  inputURL,
	}
	op, err := client.ImportEntities(ctx, req)
	if err != nil {
		return fmt.Errorf("ImportEntities: %w", err)
	}
	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait: %w", err)
	}
	fmt.Fprintf(w, "Entities were imported\n")
	return nil
}

// [END datastore_admin_entities_import]
