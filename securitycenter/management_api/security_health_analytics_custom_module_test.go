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

import (
	"bytes"
	"context"
	"fmt"
	// "log"
	"os"
	"strings"
	"testing"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	expr "google.golang.org/genproto/googleapis/type/expr"
)

var orgID = ""
var createdCustomModuleID = ""

// setup initializes variables in this file with entityNames to
// use for testing.
func setup(t *testing.T) string {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")

	if orgID == "" {
		t.Fatalf("GCLOUD_ORGANIZATION environment variable is not set.")
	}

	return orgID
}

// func TestMain(m *testing.M) {
// 	// Perform cleanup before running tests
// 	err := cleanupExistingCustomModules(orgID)
// 	if err != nil {
// 		log.Fatalf("Error cleaning up existing custom modules: %v", err)
// 	}

// 	// Run the tests
// 	code := m.Run()

// 	// Exit with the appropriate code
// 	os.Exit(code)
// }


// extractCustomModuleID extracts the custom module ID from the full name
func extractCustomModuleID(customModuleFullName string) string {
	trimmedFullName := strings.TrimSpace(customModuleFullName)
	parts := strings.Split(trimmedFullName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// addCustomModule creates a custom module for testing purposes
func addCustomModule(t *testing.T) (string, error) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)	

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

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
			Description: "Sample custom module for testing purpose. Please do not delete.", // Replace with the desired description.
			Predicate: &expr.Expr{
				Expression:  "has(resource.rotationPeriod) && (resource.rotationPeriod > duration('2592000s'))",
				Title:       "GCE Instance High Severity",
				Description: "Custom module to detect high severity issues on GCE instances.",
			},
			Recommendation: "Ensure proper security configurations on GCE instances.",
			ResourceSelector: &securitycentermanagementpb.CustomConfig_ResourceSelector{
				ResourceTypes: []string{"cloudkms.googleapis.com/CryptoKey"},
			},
			Severity:    securitycentermanagementpb.CustomConfig_CRITICAL,
		},
		DisplayName: "go_sample_custom_module_test", // Replace with desired Display Name.
		EnablementState: securitycentermanagementpb.SecurityHealthAnalyticsCustomModule_ENABLED,
	}

	req := &securitycentermanagementpb.CreateSecurityHealthAnalyticsCustomModuleRequest{
		Parent:                    parent,
		SecurityHealthAnalyticsCustomModule: customModule,
	}

	module, err := client.CreateSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(buf, "Created SecurityHealthAnalyticsCustomModule: %s\n", module.Name)

	customModuleFullName := module.Name
	customModuleID := extractCustomModuleID(customModuleFullName)

	// Store the created custom module ID for later use or cleanup
	createdCustomModuleID = customModuleID
	
	return createdCustomModuleID, nil
}

func cleanupCustomModule(t *testing.T, customModuleID string) error {
	orgID := setup(t)

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycentermanagementpb.DeleteSecurityHealthAnalyticsCustomModuleRequest{
		Name: fmt.Sprintf("organizations/%s/locations/global/securityHealthAnalyticsCustomModules/%s", orgID, customModuleID),
	}

	if err := client.DeleteSecurityHealthAnalyticsCustomModule(ctx, req); err != nil {
		return fmt.Errorf("failed to delete CustomModule: %w", err)
	}

	return nil
}

// func cleanupExistingCustomModules(t *testing.T, orgID string) error {
// 	ctx := context.Background()
// 	client, err := securitycentermanagement.NewClient(ctx)
// 	if err != nil {
// 		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
// 	}
// 	defer client.Close()

// 	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

// 	// List all existing custom modules
// 	req := &securitycentermanagementpb.ListSecurityHealthAnalyticsCustomModulesRequest{
// 		Parent: parent,
// 	}

// 	it := client.ListSecurityHealthAnalyticsCustomModules(ctx, req)
// 	for {
// 		module, err := it.Next()
// 		if err != nil {
// 			if err.Error() == "iterator done" {
// 				break
// 			}
// 			return fmt.Errorf("failed to list CustomModules: %w", err)
// 		}

// 		// Check if the custom module name starts with 'go_sample_'
// 		if strings.HasPrefix(module.DisplayName, "go_sample_") {

// 			customModuleID := extractCustomModuleID(module.Name)
// 			// Delete the custom module
// 			err := cleanupCustomModule(t, customModuleID)
// 			if err != nil {
// 				return fmt.Errorf("failed to delete existing CustomModule: %w", err)
// 			}
// 			fmt.Printf("Deleted existing CustomModule: %s\n", module.Name)
// 		}
// 	}

// 	return nil
// }

// TestCreateCustomModule verifies the Create functionality
func TestCreateCustomModule(t *testing.T) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)	

	// Call Create
	err := createSecurityHealthAnalyticsCustomModule(buf, parent)

	if err != nil {
		t.Fatalf("createCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("createCustomModule() got: %s want %s", got, orgID)
	}

	// Cleanup
	cleanupCustomModule(t, createdCustomModuleID)
}

// TestGetCustomModule verifies the Get functionality
func TestGetCustomModule(t *testing.T) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	createdCustomModuleID, err := addCustomModule(t);

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err = getSecurityHealthAnalyticsCustomModule(buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("getSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("getSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}

	// Cleanup
	cleanupCustomModule(t, createdCustomModuleID)
}

// TestUpdateCustomModule verifies the Update functionality
func TestUpdateCustomModule(t *testing.T) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	createdCustomModuleID, err := addCustomModule(t);

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	// Call Update
	err = updateSecurityHealthAnalyticsCustomModule(buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("updateSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("updateCustomModule() got: %s want %s", got, orgID)
	}

	// Cleanup
	cleanupCustomModule(t, createdCustomModuleID)

}

// TestListCustomModule verifies the List functionality
func TestListCustomModule(t *testing.T) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	createdCustomModuleID, err := addCustomModule(t);

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listSecurityHealthAnalyticsCustomModule(buf, parent)

	if err != nil {
		t.Fatalf("listSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}

	// Cleanup
	cleanupCustomModule(t, createdCustomModuleID)
}
