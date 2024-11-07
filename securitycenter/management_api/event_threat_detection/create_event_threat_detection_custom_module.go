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

// [START securitycenter_management_api_create_event_threat_detection_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// createEventThreatDetectionCustomModule creates a custom module for Security Health Analytics.
func createEventThreatDetectionCustomModule(w io.Writer, parent string) error {
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

	// Define the metadata and other config parameters as a map
	configMap := map[string]interface{}{
		"metadata": map[string]interface{}{
			"severity": "MEDIUM",
			//Replace with the desired description.
			"description":    "Sample custom module for testing purpose. Please do not delete.",
			"recommendation": "na",
		},
		"ips": []interface{}{"0.0.0.0"},
	}

	// Convert the map to a Struct
	configStruct, err := structpb.NewStruct(configMap)
	if err != nil {
		return fmt.Errorf("structpb.NewStruct: %w", err)
	}

	// Define the Event Threat Detection custom module configuration
	customModule := &securitycentermanagementpb.EventThreatDetectionCustomModule{
		Config: configStruct,
		//Replace with desired Display Name.
		DisplayName:     "go_sample_etd_custom_module",
		EnablementState: securitycentermanagementpb.EventThreatDetectionCustomModule_ENABLED,
		Type:            "CONFIGURABLE_BAD_IP",
	}

	req := &securitycentermanagementpb.CreateEventThreatDetectionCustomModuleRequest{
		Parent:                           parent,
		EventThreatDetectionCustomModule: customModule,
	}

	module, err := client.CreateEventThreatDetectionCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create EventThreatDetectionCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Created EventThreatDetectionCustomModule: %s\n", module.Name)
	return nil
}

// [END securitycenter_management_api_create_event_threat_detection_custom_module]
