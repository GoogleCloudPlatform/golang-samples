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

package snippets

// [START healthcare_create_dataset]
import (
	"context"
	"fmt"
	"io"
	"time"

	healthcare "google.golang.org/api/healthcare/v1"
)

// createDataset creates a dataset.
func createDataset(w io.Writer, projectID, location, datasetID string) error {
	// Set a deadline for the dataset to become initialized.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	datasetsService := healthcareService.Projects.Locations.Datasets

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, location)

	resp, err := datasetsService.Create(parent, &healthcare.Dataset{}).DatasetId(datasetID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	// The dataset is not always ready to use immediately, instead a long-running operation is returned.
	// This is how you might poll the operation to ensure the dataset is fully initialized before proceeding.
	// Initialization usually takes less than a minute.
	for !resp.Done {
		time.Sleep(15 * time.Second)
		resp, err = datasetsService.Operations.Get(resp.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("Operations.Get(%s): %w", resp.Name, err)
		}
	}

	fmt.Fprintf(w, "Created dataset: %q\n", resp.Name)
	return nil
}

// [END healthcare_create_dataset]
