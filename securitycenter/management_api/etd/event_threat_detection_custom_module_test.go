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

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	iterator "google.golang.org/api/iterator"
	// expr "google.golang.org/genproto/googleapis/type/expr"
)

var orgID = ""
var createdCustomModuleID = ""

func TestMain(m *testing.M) {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")

	if orgID == "" {
		log.Fatalf("GCLOUD_ORGANIZATION environment variable is not set.")
	}

	// Perform cleanup before running tests
	err := cleanupExistingCustomModules(orgID)
	if err != nil {
		log.Fatalf("Error cleaning up existing custom modules: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
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
	// Create unique display name
	rand.Seed(time.Now().UnixNano())
	displayName := fmt.Sprintf("go_sample_etd_custom_module_test_%d", rand.Int())

	// Define the Event Threat Detection custom module configuration
	customModule := &securitycentermanagementpb.EventThreatDetectionCustomModule{
		Config: &securitycentermanagementpb.EventThreatDetectionCustomModule_Config{
			Metadata: &securitycentermanagementpb.EventThreatDetectionCustomModule_Config_Metadata{
				//Replace with the desired severity.
				Severity:      "MEDIUM",
				//Replace with the desired description.
				Description:   "Sample custom module for testing purpose. Please do not delete.",
				Recommendation: "na",
			},
			Ips: []string{"0.0.0.0"},
		},
		//Replace with desired Display Name.
		DisplayName:    displayName,
		EnablementState: securitycentermanagementpb.EventThreatDetectionCustomModule_ENABLED,
		Type:           securitycentermanagementpb.EventThreatDetectionCustomModule_CONFIGURABLE_BAD_IP,
	}

	req := &securitycentermanagementpb.CreateEventThreatDetectionCustomModuleRequest{
		Parent:                      parent,
		EventThreatDetectionCustomModule: customModule,
	}

	module, err := client.CreateEventThreatDetectionCustomModule(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create EventThreatDetectionCustomModule: %w", err)
	}		

	fmt.Fprintf(&buf, "Created EventThreatDetectionCustomModule: %s\n", module.Name)

	customModuleFullName := module.Name
	customModuleID := extractCustomModuleID(customModuleFullName)

	// Store the created custom module ID for later use or cleanup
	createdCustomModuleID = customModuleID

	return createdCustomModuleID, nil
}

func cleanupCustomModule(customModuleID string) error {

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycentermanagement.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycentermanagementpb.DeleteEventThreatDetectionCustomModuleRequest{
		Name: fmt.Sprintf("organizations/%s/locations/global/eventThreatDetectionCustomModules/%s", orgID, customModuleID),
	}

	if err := client.DeleteEventThreatDetectionCustomModule(ctx, req); err != nil {
		return fmt.Errorf("failed to delete CustomModule: %w", err)
	}

	return nil
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
	req := &securitycentermanagementpb.ListEventThreatDetectionCustomModulesRequest{
		Parent: parent,
	}

	it := client.ListEventThreatDetectionCustomModules(ctx, req)
	for {
		module, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to list CustomModules: %w", err)
		}

		// Check if the custom module name starts with 'go_sample_etd_custom'
		if strings.HasPrefix(module.DisplayName, "go_sample_etd_custom") {

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

// TestCreateEtdCustomModule verifies the Create functionality
func TestCreateEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Create
	err := createEventThreatDetectionCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("createCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("createCustomModule() got: %s want %s", got, orgID)
	}
}

// TestGetCustomModule verifies the Get functionality
func TestGetCustomModule(t *testing.T) {
	var buf bytes.Buffer

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err = getEventThreatDetectionCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("getEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("getEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
	}
}
