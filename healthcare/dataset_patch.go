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

// [START healthcare_patch_dataset]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// patchDataset updates (patches) a dataset by updating its timezone..
func patchDataset(w io.Writer, projectID, location, datasetID, newTimeZone string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	datasetsService := healthcareService.Projects.Locations.Datasets

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)

	if _, err := datasetsService.Patch(name, &healthcare.Dataset{
		TimeZone: newTimeZone,
	}).UpdateMask("timeZone").Do(); err != nil {
		return fmt.Errorf("Patch: %v", err)
	}

	fmt.Fprintf(w, "Patched dataset %s with timeZone %s\n", datasetID, newTimeZone)

	return nil
}

// [END healthcare_patch_dataset]
