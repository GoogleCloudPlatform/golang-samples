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
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	iterator "google.golang.org/api/iterator"
	expr "google.golang.org/genproto/googleapis/type/expr"
)

var orgID = ""
var createdCustomModuleID = ""
var mu sync.Mutex

func TestMain(m *testing.M) {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")

	if orgID == "" {
		log.Fatalf("GCLOUD_ORGANIZATION environment variable is not set.")
	}

	// Perform cleanup before running tests
	if err := retryOperation(func() error {
		return cleanupExistingCustomModules(orgID)
	}, 5, 2*time.Second); err != nil {
		log.Fatalf("Error cleaning up existing custom modules: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

func retryOperation(operation func() error, retries int, baseDelay time.Duration) error {
	for i := 0; i <= retries; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		if i < retries {
			delay := time.Duration(math.Pow(2, float64(i))) * baseDelay
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("operation failed after %d retries", retries)
}

func cleanupExistingCustomModules(orgID string) error {
	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// List all existing custom modules
	req := &securitycentermanagementpb.ListSecurityHealthAnalyticsCustomModulesRequest{
		Parent: parent,
	}

	it := client.ListSecurityHealthAnalyticsCustomModules(ctx, req)
	for {
		module, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to list CustomModules: %w", err)
		}

		// Check if the custom module name starts with 'go_sample_sha_custom'
		if strings.HasPrefix(module.DisplayName, "go_sample_sha_custom") {

			customModuleID := extractCustomModuleID(module.Name)
			// Delete the custom module
			err := cleanupCustomModule(customModuleID)
			if err != nil {
				return fmt.Errorf("failed to delete existing CustomModule: %w", err)
			}
			fmt.Printf("Deleted existing CustomModule: %s\n", module.Name)
		}
	}

	return nil
}

func cleanupCustomModule(customModuleID string) error {

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
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Generate a unique suffix
	uniqueSuffix := fmt.Sprintf("%d_%d", time.Now().Unix(), rand.Intn(1000))
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

// extractCustomModuleID extracts the custom module ID from the full name
func extractCustomModuleID(customModuleFullName string) string {
	trimmedFullName := strings.TrimSpace(customModuleFullName)
	parts := strings.Split(trimmedFullName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// TestDeleteCustomModule verifies the List functionality
func TestDeleteCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

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

// TestCreateSHACustomModule verifies the Create functionality
func TestCreateSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Create
	err := createSecurityHealthAnalyticsCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("createCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("createCustomModule() got: %s want %s", got, orgID)
	}
}

// TestListDescendantSHACustomModule verifies the List Descendant functionality
func TestListDescendantSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestListDescendantSHACustomModule: %v", err)
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

// TestGetSHACustomModule verifies the Get functionality
func TestGetSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestGetSHACustomModule: %v", err)
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

// TestSimulateSHACustomModule verifies the Create functionality
func TestSimulateSHACustomModule(t *testing.T) {
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

// TestListEffectiveSHACustomModule verifies the List Effective functionality
func TestListEffectiveSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestListEffectiveSHACustomModule: %v", err)
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

// TestUpdateSHACustomModule verifies the Update functionality
func TestUpdateSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestUpdateSHACustomModule: %v", err)
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

// TestGetEffectiveSHACustomModule verifies the Get Effective functionality
func TestGetEffectiveSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestGetEffectiveSHACustomModule: %v", err)
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

// TestListSHACustomModule verifies the List functionality
func TestListSHACustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestListSHACustomModule: %v", err)
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
