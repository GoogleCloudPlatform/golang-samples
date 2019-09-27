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

// [START healthcare_list_hl7v2_messages]
import (
	"context"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// listHL7V2Messages prints a list of HL7V2 messages to w.
func listHL7V2Messages(w io.Writer, projectID, location, datasetID, hl7V2StoreID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	messagesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores.Messages

	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7v2Stores/%s", projectID, location, datasetID, hl7V2StoreID)

	resp, err := messagesService.List(parent).Do()
	if err != nil {
		return fmt.Errorf("Create: %v", err)
	}

	fmt.Fprintln(w, "HL7V2 messages:")
	for _, s := range resp.Messages {
		fmt.Fprintln(w, s)
	}
	return nil
}

// [END healthcare_list_hl7v2_messages]
