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

// [START healthcare_list_datasets]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// listDatasets prints a list of datasets to w.
func listDatasets(w io.Writer, projectID string, location string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	datasetsService := healthcareService.Projects.Locations.Datasets

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, location)

	resp, err := datasetsService.List(parent).Do()
	if err != nil {
		return fmt.Errorf("List: %w", err)
	}

	fmt.Fprintln(w, "Datasets:")
	for _, d := range resp.Datasets {
		fmt.Fprintln(w, d.Name)
	}

	return nil
}

// [END healthcare_list_datasets]
