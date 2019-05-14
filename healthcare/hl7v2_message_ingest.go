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

// [START healthcare_ingest_hl7v2_message]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	healthcare "google.golang.org/api/healthcare/v1beta1"
)

// ingestHL7V2Message ingests an HL7V2 message.
func ingestHL7V2Message(w io.Writer, projectID, location, datasetID, hl7V2StoreID, messageFile string) error {
	ctx := context.Background()

	hl7v2message, err := ioutil.ReadFile(messageFile)
	if err != nil {
		return fmt.Errorf("ReadFile: %v", err)
	}

	healthcareService, err := healthcare.NewService(ctx)
	if err != nil {
		return fmt.Errorf("healthcare.NewService: %v", err)
	}

	messagesService := healthcareService.Projects.Locations.Datasets.Hl7V2Stores.Messages

	req := &healthcare.IngestMessageRequest{
		Message: &healthcare.Message{
			Data: base64.StdEncoding.EncodeToString(hl7v2message),
		},
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7V2Stores/%s", projectID, location, datasetID, hl7V2StoreID)

	resp, err := messagesService.Ingest(parent, req).Do()
	if err != nil {
		return fmt.Errorf("Create: %v", err)
	}

	fmt.Fprintf(w, "Ingested HL7V2 message: %q\n", resp.Message.Name)
	return nil
}

// [END healthcare_ingest_hl7v2_message]
