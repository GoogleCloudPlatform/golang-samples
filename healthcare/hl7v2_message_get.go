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

// [START healthcare_get_hl7v2_message]
import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	healthcare "google.golang.org/api/healthcare/v1"
)

// getHL7V2Message gets an HL7V2 message.
func getHL7V2Message(w io.Writer, projectID, location, datasetID, hl7V2StoreID, hl7V2MessageID string) error {
	ctx := context.Background()

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	messagesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores.Messages

	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7V2Stores/%s/messages/%s", projectID, location, datasetID, hl7V2StoreID, hl7V2MessageID)
	message, err := messagesService.Get(name).Do()
	if err != nil {
		return fmt.Errorf("Get: %w", err)
	}

	rawData, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		return fmt.Errorf("base64.DecodeString: %w", err)
	}

	fmt.Fprintf(w, "Got HL7V2 message.\n")
	fmt.Fprintf(w, "Raw length: %d.\n", len(rawData))
	fmt.Fprintf(w, "Parsed data:\n")
	parsedJSON, _ := json.MarshalIndent(message.ParsedData, "", "  ")
	fmt.Fprintf(w, "%s", parsedJSON)
	return nil
}

// [END healthcare_get_hl7v2_message]
