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

// [START securitycenter_management_api_create_security_health_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycenterpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	exprpb "google.golang.org/genproto/googleapis/type/expr"
)

// CreateSecurityHealthAnalyticsCustomModule creates a custom module for Security Health Analytics.
func CreateSecurityHealthAnalyticsCustomModule(w io.Writer, parent string, customModuleID string) error {
	// parent: Use any one of the following options:
	// - organizations/{organization_id}/locations/{location_id}
	// - folders/{folder_id}/locations/{location_id}
	// - projects/{project_id}/locations/{location_id}

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	// Define the custom module configuration
	customModule := &securitycenterpb.SecurityHealthAnalyticsCustomModule{
		Name: fmt.Sprintf("%s/securityHealthAnalyticsCustomModules/%s", parent, customModuleID),
		DisplayName: "Custom Module for testing",
		EnablementState: securitycenterpb.SecurityHealthAnalyticsCustomModule_ENABLED,
		CustomConfig: &securitycenterpb.CustomConfig{
			Predicate: &exprpb.Expr{
				Expression:  "resource.type == \"gce_instance\" && severity == \"HIGH\"",
				Title:       "GCE Instance High Severity",
				Description: "Custom module to detect high severity issues on GCE instances.",
			},
			Description: "A custom module for detecting high severity issues on GCE instances.",
			Severity:    securitycenterpb.CustomConfig_CRITICAL,
			ResourceSelector: &securitycenterpb.CustomConfig_ResourceSelector{
				ResourceTypes: []string{"CryptoKey"},
			},
			Recommendation: "Ensure proper security configurations on GCE instances.",
			CustomOutput: &securitycenterpb.CustomConfig_CustomOutputSpec{
				Properties: []*securitycenterpb.CustomConfig_CustomOutputSpec_Property{
					{
						Name: "example_property",
						ValueExpression: &exprpb.Expr{
							Expression: "resource.name",
						},
					},
				},
			},
		},
	}

	req := &securitycenterpb.CreateSecurityHealthAnalyticsCustomModuleRequest{
		Parent:                    parent,
		SecurityHealthAnalyticsCustomModule: customModule,
	}

	module, err := client.CreateSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Created SecurityHealthAnalyticsCustomModule: %s\n", module.Name)
	return nil
}

// [END securitycenter_management_api_create_security_health_custom_module]
