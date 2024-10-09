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

// [START securitycenter_management_api_get_security_health_analytics_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycenterpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
)

// GetSecurityHealthAnalyticsCustomModule retrieves a specific custom module by its name.
func getSecurityHealthAnalyticsCustomModule(w io.Writer, parent string, customModuleID string) error {
	// parent: Use any one of the following options:
	//             - organizations/{organization_id}/locations/{location_id}
	//             - folders/{folder_id}/locations/{location_id}
	//             - projects/{project_id}/locations/{location_id}
	// customModuleID := "your-module-id"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.GetSecurityHealthAnalyticsCustomModuleRequest{
		Name: fmt.Sprintf("%s/securityHealthAnalyticsCustomModules/%s", parent, customModuleID),
	}

	module, err := client.GetSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to get SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Retrieved SecurityHealthAnalyticsCustomModule: %s\n", module.Name)
	return nil
}

// [END securitycenter_get_security_health_analytics_custom_module_v1]
