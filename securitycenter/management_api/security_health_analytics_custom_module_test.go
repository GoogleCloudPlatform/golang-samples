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
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	"github.com/google/uuid"
	expr "google.golang.org/genproto/googleapis/type/expr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var orgID = ""
var createdCustomModuleID = ""
var createdModules []string

func TestMain(m *testing.M) {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")

	if orgID == "" {
		log.Fatalf("GCLOUD_ORGANIZATION environment variable is not set.")
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

// AddModuleToCleanup registers a module for cleanup.
func AddModuleToCleanup(moduleID string) {
	createdModules = append(createdModules, moduleID)
}

// PrintAllCreatedModules prints all created custom modules.
func PrintAllCreatedModules() {

	if len(createdModules) == 0 {
		fmt.Println("No custom modules were created.")
	} else {
		fmt.Println("Created Custom Modules:")
		for _, module := range createdModules {
			fmt.Println(module)
		}
	}
}

// CleanupCreatedModules deletes all created custom modules.
func CleanupCreatedModules() {

	if len(createdModules) == 0 {
		fmt.Println("No custom modules to clean up.")
		return
	}

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create SecurityCenter client: %v", err)
	}
	defer client.Close()

	for len(createdModules) > 0 {
		moduleID := createdModules[0]
		if !CustomModuleExists(moduleID) {
			fmt.Printf("Module not found (already deleted): %s\n", moduleID)
			createdModules = createdModules[1:]
			continue
		}
		err := client.DeleteSecurityHealthAnalyticsCustomModule(ctx, &securitycentermanagementpb.DeleteSecurityHealthAnalyticsCustomModuleRequest{
			Name: fmt.Sprintf("organizations/%s/locations/global/securityHealthAnalyticsCustomModules/%s", orgID, moduleID),
		})

		if err != nil {
			fmt.Printf("Failed to delete module %s: %v\n", moduleID, err)
			return
		}
		fmt.Printf("Deleted custom module: %s\n", moduleID)
		createdModules = createdModules[1:]
	}
}

// CustomModuleExists checks if a module exists.
func CustomModuleExists(moduleID string) bool {
	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	_, err = client.GetSecurityHealthAnalyticsCustomModule(ctx, &securitycentermanagementpb.GetSecurityHealthAnalyticsCustomModuleRequest{
		Name: fmt.Sprintf("organizations/%s/locations/global/securityHealthAnalyticsCustomModules/%s", orgID, moduleID),
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return false
		}
		log.Printf("Error checking module existence: %v", err)
	}
	return true
}

// CleanupAfterTests is a helper for test cleanup.
func CleanupAfterTests(t *testing.T) {
	t.Cleanup(func() {
		PrintAllCreatedModules()
		fmt.Println("Cleaning up created custom modules...")
		CleanupCreatedModules()
	})
}

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
func addCustomModule() (string, error) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()
	uniqueSuffix := uuid.New().String()

	// Remove invalid characters (anything that isn't alphanumeric or an underscore)
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	uniqueSuffix = re.ReplaceAllString(uniqueSuffix, "_")

	// Create unique display name
	displayName := fmt.Sprintf("go_sample_sha_custom_module_test_%s", uniqueSuffix)

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
			Severity: securitycentermanagementpb.CustomConfig_CRITICAL,
		},
		// Replace with desired Display Name.
		DisplayName:     displayName,
		EnablementState: securitycentermanagementpb.SecurityHealthAnalyticsCustomModule_ENABLED,
	}

	req := &securitycentermanagementpb.CreateSecurityHealthAnalyticsCustomModuleRequest{
		Parent:                              parent,
		SecurityHealthAnalyticsCustomModule: customModule,
	}

	module, err := client.CreateSecurityHealthAnalyticsCustomModule(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create SecurityHealthAnalyticsCustomModule: %w", err)
	}

	fmt.Fprintf(&buf, "Created SecurityHealthAnalyticsCustomModule: %s\n", module.Name)

	customModuleFullName := module.Name
	customModuleID := extractCustomModuleID(customModuleFullName)

	// Store the created custom module ID for later use or cleanup
	createdCustomModuleID = customModuleID

	return createdCustomModuleID, nil
}

// TestDeleteCustomModule verifies the List functionality
func TestDeleteCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = deleteSecurityHealthAnalyticsCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("deleteSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, createdCustomModuleID) {
		t.Fatalf("deleteSecurityHealthAnalyticsCustomModule() got: %s want %s", got, createdCustomModuleID)
	}
}

// TestCreateCustomModule verifies the Create functionality
func TestCreateCustomModule(t *testing.T) {
	var buf bytes.Buffer
	var createModulePath = ""

	CleanupAfterTests(t)

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Create
	err := createSecurityHealthAnalyticsCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("createCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if got == "" {
		t.Errorf("createSecurityHealthAnalyticsCustomModule() returned an empty string")
		return
	}

	fmt.Printf("Response: %v\n", got)

	parts := strings.Split(got, ":")
	if len(parts) > 0 {
		createModulePath = parts[len(parts)-1]
	}

	AddModuleToCleanup(extractCustomModuleID(createModulePath))

	if !strings.Contains(got, orgID) {
		t.Fatalf("createCustomModule() got: %s want %s", got, orgID)
	}
}

// TestListDescendantCustomModule verifies the List Descendant functionality
func TestListDescendantCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listDescendantSecurityHealthAnalyticsCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("listDescendantSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listDescendantSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}
}

// TestGetCustomModule verifies the Get functionality
func TestGetCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err = getSecurityHealthAnalyticsCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("getSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("getSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}
}

// TestSimulateCustomModule verifies the Create functionality
func TestSimulateCustomModule(t *testing.T) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Simulate
	err := simulateSecurityHealthAnalyticsCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("simulateCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if want := "no_violation"; !strings.Contains(got, want) {
		t.Fatalf("simulateCustomModule() got: %s want %s", got, want)
	}
}

// TestListEffectiveCustomModule verifies the List Effective functionality
func TestListEffectiveCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listEffectiveSecurityHealthAnalyticsCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("listEffectiveSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listEffectiveSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}
}

// TestUpdateCustomModule verifies the Update functionality
func TestUpdateCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	// Call Update
	err = updateSecurityHealthAnalyticsCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("updateSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("updateCustomModule() got: %s want %s", got, orgID)
	}
}

// TestGetEffectiveCustomModule verifies the Get Effective functionality
func TestGetEffectiveCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err = getEffectiveSecurityHealthAnalyticsCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("getEffectiveSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("getEffectiveSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}
}

// TestListCustomModule verifies the List functionality
func TestListCustomModule(t *testing.T) {
	var buf bytes.Buffer

	CleanupAfterTests(t)

	createdCustomModuleID, err := addCustomModule()
	AddModuleToCleanup(createdCustomModuleID)

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listSecurityHealthAnalyticsCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("listSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listSecurityHealthAnalyticsCustomModule() got: %s want %s", got, orgID)
	}
}
