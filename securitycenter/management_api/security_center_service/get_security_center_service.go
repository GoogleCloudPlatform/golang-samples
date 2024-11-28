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

package security_center_service

// [START securitycenter_get_security_center_service]

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
)

// getSecurityCenterService retrieves a specific Security Center service by its name.
func getSecurityCenterService(w io.Writer, parent string, service string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}/locations/{location_id}
	//             - folders/{folder_id}/locations/{location_id}
	//             - projects/{project_id}/locations/{location_id}
	// Replace `{service}` with one of the valid values:
	// container-threat-detection, event-threat-detection, security-health-analytics, vm-threat-detection, web-security-scanner
	// service := "security-center-service-name"
	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycentermanagementpb.GetSecurityCenterServiceRequest{
		Name: fmt.Sprintf("%s/securityCenterServices/%s", parent, service),
	}

	response, err := client.GetSecurityCenterService(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to get SecurityCenterService: %w", err)
	}

	fmt.Fprintf(w, "Retrieved SecurityCenterService: %s\n", response.Name)
	return nil
}

// [END securitycenter_get_security_center_service]
