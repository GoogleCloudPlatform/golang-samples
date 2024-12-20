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

package event_threat_detection

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
	"google.golang.org/protobuf/types/known/structpb"
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
	uniqueSuffix := fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Intn(1000))
	// Create unique display name
	displayName := fmt.Sprintf("go_sample_etd_custom_module_test_%s", uniqueSuffix)

	// Define the metadata and other config parameters as a map
	configMap := map[string]interface{}{
		"metadata": map[string]interface{}{
			"severity": "MEDIUM",
			//Replace with the desired description.
			"description":    "Sample custom module for testing purpose. Please do not delete.",
			"recommendation": "na",
		},
		"ips": []interface{}{"0.0.0.0"},
	}

	// Convert the map to a Struct
	configStruct, err := structpb.NewStruct(configMap)
	if err != nil {
		return "", fmt.Errorf("structpb.NewStruct: %w", err)
	}

	// Define the Event Threat Detection custom module configuration
	customModule := &securitycentermanagementpb.EventThreatDetectionCustomModule{
		Config: configStruct,
		//Replace with desired Display Name.
		DisplayName:     displayName,
		EnablementState: securitycentermanagementpb.EventThreatDetectionCustomModule_ENABLED,
		Type:            "CONFIGURABLE_BAD_IP",
	}

	req := &securitycentermanagementpb.CreateEventThreatDetectionCustomModuleRequest{
		Parent:                           parent,
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

// extractCustomModuleID extracts the custom module ID from the full name
func extractCustomModuleID(customModuleFullName string) string {
	trimmedFullName := strings.TrimSpace(customModuleFullName)
	parts := strings.Split(trimmedFullName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
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
func TestGetEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestGetEtdCustomModule: %v", err)
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

// TestUpdateCustomModule verifies the Update functionality
func TestUpdateEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestUpdateEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	// Call Update
	err = updateEventThreatDetectionCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("updateEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("updateCustomModule() got: %s want %s", got, orgID)
	}
}

// TestDeleteCustomModule verifies the List functionality
func TestDeleteEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestDeleteEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = deleteEventThreatDetectionCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("deleteEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, createdCustomModuleID) {
		t.Fatalf("deleteEventThreatDetectionCustomModule() got: %s want %s", got, createdCustomModuleID)
	}
}

// TestListEtdCustomModule verifies the List functionality
func TestListEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestListEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listEventThreatDetectionCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("listEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
	}
}

// TestListEffectiveEtdCustomModule verifies the List functionality
func TestListEffectiveEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestListEffectiveEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listEffectiveEventThreatDetectionCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("listEffectiveEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listEffectiveEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
	}
}

// TestGetEffectiveEtdCustomModule verifies the Get functionality
func TestGetEffectiveEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	createdCustomModuleID, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestGetEffectiveEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err = getEffectiveEventThreatDetectionCustomModule(&buf, parent, createdCustomModuleID)

	if err != nil {
		t.Fatalf("getEffectiveEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("getEffectiveEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
	}
}

// TestListDescendantEtdCustomModule verifies the List functionality
func TestListDescendantEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestListDescendantEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = listDescendantEventThreatDetectionCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("listDescendantEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listDescendantEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
	}
}

// TestValidateEtdCustomModule verifies the List functionality
func TestValidateEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	mu.Lock()
	defer mu.Unlock()

	_, err := addCustomModule()

	if err != nil {
		t.Fatalf("Could not setup test environment at TestValidateEtdCustomModule: %v", err)
		return
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err = validateEventThreatDetectionCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("validateEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	// Check that the response indicates successful validation
	expectedMessage := "Validation successful: No errors found."

	if !strings.Contains(got, expectedMessage) {
		t.Fatalf("validateEventThreatDetectionCustomModule() got: %s want %s", got, expectedMessage)
	}
}
