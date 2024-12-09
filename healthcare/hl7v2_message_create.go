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

// [START healthcare_create_hl7v2_message]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	healthcare "google.golang.org/api/healthcare/v1"
)

// createHL7V2Message creates an HL7V2 message.
func createHL7V2Message(w io.Writer, projectID, location, datasetID, hl7V2StoreID, messageFile string) error {
	ctx := context.Background()

	hl7v2message, err := os.ReadFile(messageFile)
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %w", err)
	}

	messagesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores.Messages

	req := &healthcare.CreateMessageRequest{
		Message: &healthcare.Message{
			Data: base64.StdEncoding.EncodeToString(hl7v2message),
		},
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7V2Stores/%s", projectID, location, datasetID, hl7V2StoreID)
	resp, err := messagesService.Create(parent, req).Do()
	if err != nil {
		return fmt.Errorf("messagesService.Create: %w", err)
	}

	fmt.Fprintf(w, "Created HL7V2 message: %q\n", resp.Name)
	return nil
}

// [END healthcare_create_hl7v2_message]
