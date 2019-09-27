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

// [START healthcare_export_fhir_resources_gcs]
import (
	"context"
	"fmt"
	"io"
	"time"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// exportFHIRResource exports the resources in the FHIR store.
func exportFHIRResource(w io.Writer, projectID, location, datasetID, fhirStoreID, gcsURIPrefix string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.FhirStores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)
	req := &healthcare.ExportResourcesRequest{
		GcsDestination: &healthcare.GoogleCloudHealthcareV1beta1FhirRestGcsDestination{
			UriPrefix: gcsURIPrefix,
		},
	}

	op, err := storesService.Export(name, req).Do()
	if err != nil {
		return fmt.Errorf("Export: %v", err)
	}

	operationsService := healthcareService.Projects.Locations.Datasets.Operations
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			newOp, err := operationsService.Get(op.Name).Do()
			if err != nil {
				return fmt.Errorf("operationsService.Get(%q): %v", op.Name, err)
			}
			if newOp.Done {
				if newOp.Error != nil {
					return fmt.Errorf("export operation %q completed with error: %v", op.Name, newOp.Error)
				}
				return nil
			}
		}
	}
}

// [END healthcare_export_fhir_resources_gcs]
