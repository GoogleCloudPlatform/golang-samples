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

// [START healthcare_patch_hl7v2_store]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// patchHL7V2Store updates (patches) a HL7V2 store by updating its Pub/sub topic name.
func patchHL7V2Store(w io.Writer, projectID, location, datasetID, hl7v2StoreID, topicName string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	storesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7v2Stores/%s", projectID, location, datasetID, hl7v2StoreID)

	if _, err := storesService.Patch(name, &healthcare.Hl7V2Store{
		NotificationConfig: &healthcare.NotificationConfig{
			PubsubTopic: topicName, // format is "projects/*/locations/*/topics/*"
		},
	}).UpdateMask("notificationConfig").Do(); err != nil {
		return fmt.Errorf("Patch: %v", err)
	}

	fmt.Fprintf(w, "Patched HL7V2 store %s with Pub/sub topic %s\n", datasetID, topicName)

	return nil
}

// [END healthcare_patch_hl7v2_store]
