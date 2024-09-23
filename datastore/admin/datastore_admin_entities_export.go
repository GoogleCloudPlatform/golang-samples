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

// [START datastore_admin_entities_export]
import (
	"context"
	"fmt"
	"io"

	admin "cloud.google.com/go/datastore/admin/apiv1"
	"cloud.google.com/go/datastore/admin/apiv1/adminpb"
)

// entitiesExport exports a copy of all or a subset of entities from
// Datastore to another storage system, such as Cloud Storage.
func entitiesExport(w io.Writer, projectID, outputURLPrefix string) (*adminpb.ExportEntitiesResponse, error) {
	// projectID := "project-id"
	// outputURLPrefix := "gs://bucket-name"
	ctx := context.Background()
	client, err := admin.NewDatastoreAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("admin.NewDatastoreAdminClient: %w", err)
	}
	defer client.Close()

	req := &adminpb.ExportEntitiesRequest{
		ProjectId:       projectID,
		OutputUrlPrefix: outputURLPrefix,
	}
	op, err := client.ExportEntities(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ExportEntities: %w", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("Wait: %w", err)
	}
	fmt.Fprintln(w, "Entities were exported")
	return resp, nil
}

// [END datastore_admin_entities_export]
