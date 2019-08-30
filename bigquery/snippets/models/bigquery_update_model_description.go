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

package models

// [START bigquery_update_model_description]

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

func updateModelDescription(projectID, datasetID, modelID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// modelID := "mymodel"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}

	model := client.Dataset(datasetID).Model(modelID)
	oldMeta, err := model.Metadata(ctx)
	if err != nil {
		return fmt.Errorf("Metadata: %v", err)
	}
	update := bigquery.ModelMetadataToUpdate{
		Description: "This model was modified from a Go program",
	}
	if _, err = model.Update(ctx, update, oldMeta.ETag); err != nil {
		return fmt.Errorf("Update: %v", err)
	}
	return nil
}

// [END bigquery_update_model_description]
