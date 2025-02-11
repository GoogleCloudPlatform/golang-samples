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
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	securitycentermanagement "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycentermanagementpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

var orgID = ""
var createdCustomModuleID = ""
var moduleID = ""
var sharedModules []string

func TestMain(m *testing.M) {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")

	if orgID == "" {
		log.Fatalf("GCLOUD_ORGANIZATION environment variable is not set.")
	}

	setupSharedModules()

	// Run the tests
	code := m.Run()

	PrintAllCreatedModules()
	CleanupSharedModules()

	// Exit with the appropriate code
	os.Exit(code)
}

func setupSharedModules() {
	for i := 0; i < 3; i++ {
		moduleID, _ = addCustomModule()
		if moduleID != "" {
			sharedModules = append(sharedModules, moduleID)
		}
	}
}

// AddModuleToCleanup registers a module for cleanup.
func AddModuleToCleanup(moduleID string) {
	sharedModules = append(sharedModules, moduleID)
}

// PrintAllCreatedModules prints all created custom modules.
func PrintAllCreatedModules() {

	if len(sharedModules) == 0 {
		fmt.Println("No custom modules were created.")
	} else {
		fmt.Println("Created Custom Modules:")
		for _, module := range sharedModules {
			fmt.Println(module)
		}
	}
}

// CleanupSharedModules deletes all created custom modules.
func CleanupSharedModules() {

	if len(sharedModules) == 0 {
		fmt.Println("No custom modules to clean up.")
		return
	}

	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create SecurityCenter client: %v", err)
	}
	defer client.Close()

	for len(sharedModules) > 0 {
		moduleID = sharedModules[0]
		if !CustomModuleExists(moduleID) {
			fmt.Printf("Module not found (already deleted): %s\n", moduleID)
			sharedModules = sharedModules[1:]
			continue
		}
		err := client.DeleteEventThreatDetectionCustomModule(ctx, &securitycentermanagementpb.DeleteEventThreatDetectionCustomModuleRequest{
			Name: fmt.Sprintf("organizations/%s/locations/global/eventThreatDetectionCustomModules/%s", orgID, moduleID),
		})

		if err != nil {
			fmt.Printf("Failed to delete module %s: %v\n", moduleID, err)
			return
		}
		fmt.Printf("Deleted custom module: %s\n", moduleID)
		sharedModules = sharedModules[1:]
	}
}

func getRandomSharedModule() string {
	if len(sharedModules) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return sharedModules[rand.Intn(len(sharedModules))]
}

// CustomModuleExists checks if a module exists.
func CustomModuleExists(moduleID string) bool {
	ctx := context.Background()
	client, err := securitycentermanagement.NewClient(ctx)
	_, err = client.GetEventThreatDetectionCustomModule(ctx, &securitycentermanagementpb.GetEventThreatDetectionCustomModuleRequest{
		Name: fmt.Sprintf("organizations/%s/locations/global/eventThreatDetectionCustomModules/%s", orgID, moduleID),
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return false
		}
		log.Printf("Error checking module existence: %v", err)
	}
	return true
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

// TestCreateEtdCustomModule verifies the Create functionality
func TestCreateEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer
	var createModulePath = ""

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Create
	err := createEventThreatDetectionCustomModule(&buf, parent)

	if err != nil {
		t.Fatalf("createCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if got == "" {
		t.Errorf("createEventThreatDetectionCustomModule() returned an empty string")
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

// TestGetCustomModule verifies the Get functionality
func TestGetEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err := getEventThreatDetectionCustomModule(&buf, parent, moduleID)

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

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	// Call Update
	err := updateEventThreatDetectionCustomModule(&buf, parent, moduleID)

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

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err := deleteEventThreatDetectionCustomModule(&buf, parent, moduleID)

	if err != nil {
		t.Fatalf("deleteEventThreatDetectionCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, moduleID) {
		t.Fatalf("deleteEventThreatDetectionCustomModule() got: %s want %s", got, moduleID)
	}
}

// TestListEtdCustomModule verifies the List functionality
func TestListEtdCustomModule(t *testing.T) {
	var buf bytes.Buffer

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err := listEventThreatDetectionCustomModule(&buf, parent)

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

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err := listEffectiveEventThreatDetectionCustomModule(&buf, parent)

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

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err := getEffectiveEventThreatDetectionCustomModule(&buf, parent, moduleID)

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

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err := listDescendantEventThreatDetectionCustomModule(&buf, parent)

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

	moduleID = getRandomSharedModule()
	if moduleID == "" {
		t.Fatalf("No shared modules available")
	}

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err := validateEventThreatDetectionCustomModule(&buf, parent)

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
