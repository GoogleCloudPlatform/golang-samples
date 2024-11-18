// Copyright 2024 Google LLC
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

package event_threat_detection

// [START securitycenter_validate_event_threat_detection_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	// "google.golang.org/protobuf/types/known/structpb"
)

// validateEventThreatDetectionCustomModule validates a custom module for Event Threat Detection.
func validateEventThreatDetectionCustomModule(w io.Writer, parent string) error {
	// parent: Use any one of the following options:
	// - organizations/{organization_id}/locations/{location_id}
	// - folders/{folder_id}/locations/{location_id}
	// - projects/{project_id}/locations/{location_id}

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	// Define the raw JSON configuration for the Event Threat Detection custom module
	rawText := `{
	"ips": ["192.0.2.1"],
	"metadata": {
		"properties": {
			"someProperty": "someValue"
		},
		"severity": "MEDIUM"
	}
}`

	req := &securitycentermanagementpb.ValidateEventThreatDetectionCustomModuleRequest{
		Parent:  parent,
		RawText: rawText, // Use raw JSON as a string for validation
		Type:    "CONFIGURABLE_BAD_IP",
	}

	// Perform validation
	resp, err := client.ValidateEventThreatDetectionCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to validate EventThreatDetectionCustomModule: %w", err)
	}

	// Handle the response and output validation results
	if len(resp.Errors) > 0 {
		fmt.Fprintln(w, "Validation errors:")
		for _, e := range resp.Errors {
			fmt.Fprintf(w, "Field: %s, Description: %s\n", e.FieldPath, e.Description)
		}
	} else {
		fmt.Fprintln(w, "Validation successful: No errors found.")
	}

	return nil
}

// [END securitycenter_validate_event_threat_detection_custom_module]
