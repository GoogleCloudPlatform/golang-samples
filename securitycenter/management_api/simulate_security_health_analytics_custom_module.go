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

// [START securitycenter_management_api_simulate_security_health_analytics_custom_module]

import (
	"context"
	"fmt"
	"io"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	v1 "google.golang.org/genproto/googleapis/iam/v1"
	expr "google.golang.org/genproto/googleapis/type/expr"
	"google.golang.org/protobuf/types/known/structpb"
)

// simulateSecurityHealthAnalyticsCustomModule simulates a custom module for Security Health Analytics.
func simulateSecurityHealthAnalyticsCustomModule(w io.Writer, parent string) error {
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

	// Define the custom config to simulate configuration
	customConfig := &securitycentermanagementpb.CustomConfig{
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
		//Replace with the desired description.
		Description: "Sample custom module for testing purpose. Please do not delete.",
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
	}

	// Define the simulated resource data
	resourceData := map[string]interface{}{
		"resourceId": "test-resource-id",
		"name":       "test-resource-name",
	}

	// Convert map to *structpb.Struct
	resourceDataStruct, err := structpb.NewStruct(resourceData)
	if err != nil {
		return fmt.Errorf("structpb.NewStruct: %w", err)
	}

	// Define the simulated resource
	simulatedResource := &securitycentermanagementpb.SimulateSecurityHealthAnalyticsCustomModuleRequest_SimulatedResource{
		ResourceType: "cloudkms.googleapis.com/CryptoKey", // Replace with the correct resource type
		ResourceData: resourceDataStruct,
		IamPolicyData: &v1.Policy{
			Bindings: []*v1.Binding{
				{
					Role:    "roles/owner",
					Members: []string{"user:test-user@gmail.com"},
				},
			},
		},
	}

	req := &securitycentermanagementpb.SimulateSecurityHealthAnalyticsCustomModuleRequest{
		Parent:       parent,
		CustomConfig: customConfig,
		Resource:     simulatedResource,
	}

	response, err := client.SimulateSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to simulate SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(w, "Simulated SecurityHealthAnalyticsCustomModule: %s\n", response)
	return nil
}

// [END securitycenter_management_api_simulate_security_health_analytics_custom_module]
