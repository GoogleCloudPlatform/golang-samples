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

// [START securitycenter_management_api_create_security_health_analytics_custom_module]

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	expr "google.golang.org/genproto/googleapis/type/expr"
)

// createSecurityHealthAnalyticsCustomModule creates a custom module for Security Health Analytics.
func createSecurityHealthAnalyticsCustomModule(w io.Writer, parent string) error {
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

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Generate a unique suffix
	uniqueSuffix := fmt.Sprintf("%d_%d", time.Now().Unix(), rand.Intn(1000))
	// Create unique display name
	displayName := fmt.Sprintf("go_sample_sha_custom_module_%s", uniqueSuffix)

	// Define the custom module configuration
	customModule := &securitycentermanagementpb.SecurityHealthAnalyticsCustomModule{
		CustomConfig: &securitycentermanagementpb.CustomConfig{
			CustomOutput: &securitycentermanagementpb.CustomConfig_CustomOutputSpec{
				Properties: []*securitycentermanagementpb.CustomConfig_CustomOutputSpec_Property{
					{
						Name: "example_property",
						ValueExpression: &expr.Expr{
							Description: "The name of the instance",
							Expression:  "resource.name",
							Location:    "global",
							Title:       "Instance Name",
						},
					},
				},
			},
			Description: "Sample custom module for testing purpose. Please do not delete.", //Replace with the desired description.
			Predicate: &expr.Expr{
				Expression:  "has(resource.rotationPeriod) && (resource.rotationPeriod > duration('2592000s'))",
				Title:       "GCE Instance High Severity",
				Description: "Custom module to detect high severity issues on GCE instances.",
			},
			Recommendation: "Ensure proper security configurations on GCE instances.",
			ResourceSelector: &securitycentermanagementpb.CustomConfig_ResourceSelector{
				ResourceTypes: []string{"cloudkms.googleapis.com/CryptoKey"},
			},
			Severity: securitycentermanagementpb.CustomConfig_CRITICAL,
		},
		//Replace with desired Display Name.
		DisplayName:     displayName,
		EnablementState: securitycentermanagementpb.SecurityHealthAnalyticsCustomModule_ENABLED,
	}

	req := &securitycentermanagementpb.CreateSecurityHealthAnalyticsCustomModuleRequest{
		Parent:                              parent,
		SecurityHealthAnalyticsCustomModule: customModule,
	}

	module, err := client.CreateSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Created SecurityHealthAnalyticsCustomModule: %s\n", module.Name)
	return nil
}

// [END securitycenter_management_api_create_security_health_analytics_custom_module]
