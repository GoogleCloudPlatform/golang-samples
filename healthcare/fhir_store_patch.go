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

// [START healthcare_patch_fhir_store]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// patchFHIRStore updates (patches) a FHIR store by updating its Pub/sub topic name.
func patchFHIRStore(w io.Writer, projectID, location, datasetID, fhirStoreID, topicName string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.FhirStores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, fhirStoreID)

	// topicName format is "projects/*/locations/*/topics/*"
	notificationConfig := &healthcare.FhirNotificationConfig{
		PubsubTopic: topicName,
	}

	fhirStore := &healthcare.FhirStore{
		NotificationConfigs: []*healthcare.FhirNotificationConfig{notificationConfig},
	}

	patchRequest := storesService.Patch(name, fhirStore).UpdateMask("notificationConfigs")

	if _, err := patchRequest.Do(); err != nil {
		return fmt.Errorf("Patch: %w", err)
	}

	fmt.Fprintf(w, "Patched FHIR store %s with Pub/sub topic %s\n", datasetID, topicName)

	return nil
}

// [END healthcare_patch_fhir_store]
