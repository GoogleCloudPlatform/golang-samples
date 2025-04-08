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

// [START healthcare_delete_hl7v2_message]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// deleteHL7V2Message deletes an HL7V2 message.
func deleteHL7V2Message(w io.Writer, projectID, location, datasetID, hl7V2StoreID, hl7V2MessageID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	messagesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores.Messages

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7V2Stores/%s/messages/%s", projectID, location, datasetID, hl7V2StoreID, hl7V2MessageID)
	if _, err := messagesService.Delete(name).Do(); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	fmt.Fprintf(w, "Deleted HL7V2 message: %q\n", name)
	return nil
}

// [END healthcare_delete_hl7v2_message]
