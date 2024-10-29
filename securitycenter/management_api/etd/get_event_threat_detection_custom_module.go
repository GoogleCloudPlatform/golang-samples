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

package etd

// [START securitycenter_management_api_get_event_threat_detection_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
)

// getEventThreatDetectionCustomModule retrieves a specific custom module by its name.
func getEventThreatDetectionCustomModule(w io.Writer, parent string, customModuleID string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}/locations/{location_id}
	//             - folders/{folder_id}/locations/{location_id}
	//             - projects/{project_id}/locations/{location_id}
	// customModuleID := "your-module-id"
	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycentermanagementpb.GetEventThreatDetectionCustomModuleRequest{
		Name: fmt.Sprintf("%s/securityHealthAnalyticsCustomModules/%s", parent, customModuleID),
	}

	module, err := client.GetEventThreatDetectionCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to get EventThreatDetectionCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Retrieved EventThreatDetectionCustomModule: %s\n", module.Name)
	return nil
}

// [END securitycenter_management_api_get_event_threat_detection_custom_module]
