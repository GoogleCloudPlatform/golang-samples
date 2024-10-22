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

package management_api

// [START securitycenter_management_api_update_security_health_analytics_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateSecurityHealthAnalyticsCustomModule creates a custom module for Security Health Analytics.
func updateSecurityHealthAnalyticsCustomModule(w io.Writer, parent string, customModuleID string) error {
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

	// Define the custom module configuration
	customModule := &securitycentermanagementpb.SecurityHealthAnalyticsCustomModule{
		Name:            fmt.Sprintf("%s/securityHealthAnalyticsCustomModules/%s", parent, customModuleID),
		EnablementState: securitycentermanagementpb.SecurityHealthAnalyticsCustomModule_DISABLED,
	}

	req := &securitycentermanagementpb.UpdateSecurityHealthAnalyticsCustomModuleRequest{
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"enablement_state",
			},
		},
		SecurityHealthAnalyticsCustomModule: customModule,
	}

	module, err := client.UpdateSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Updated SecurityHealthAnalyticsCustomModule: %s\n", module.Name)
	return nil
}

// [END securitycenter_management_api_update_security_health_analytics_custom_module]
