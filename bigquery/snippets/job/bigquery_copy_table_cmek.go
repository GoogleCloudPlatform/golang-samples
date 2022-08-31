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

package job

// [START bigquery_copy_table_cmek]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// copyTableWithCMEK demonstrates creating a copy of a table and ensuring the copied data is
// protected with a customer managed encryption key.
func copyTableWithCMEK(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	srcTable := client.DatasetInProject("bigquery-public-data", "samples").Table("shakespeare")
	copier := client.Dataset(datasetID).Table(tableID).CopierFrom(srcTable)
	copier.DestinationEncryptionConfig = &bigquery.EncryptionConfig{
		// TODO: Replace this key with a key you have created in Cloud KMS.
		KMSKeyName: "projects/cloud-samples-tests/locations/us-central1/keyRings/test/cryptoKeys/test",
	}
	job, err := copier.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// [END bigquery_copy_table_cmek]
