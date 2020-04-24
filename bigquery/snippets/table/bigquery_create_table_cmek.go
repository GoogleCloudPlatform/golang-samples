// Copyright 2019 Google LLC
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

package table

// [START bigquery_create_table_cmek]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// createTableWithCMEK demonstrates creating a table protected with a customer managed encryption key.
func createTableWithCMEK(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// tableID := "mytableid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	tableRef := client.Dataset(datasetID).Table(tableID)
	meta := &bigquery.TableMetadata{
		EncryptionConfig: &bigquery.EncryptionConfig{
			// TODO: Replace this key with a key you have created in Cloud KMS.
			KMSKeyName: "projects/cloud-samples-tests/locations/us/keyRings/test/cryptoKeys/test",
		},
	}
	if err := tableRef.Create(ctx, meta); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_table_cmek]
