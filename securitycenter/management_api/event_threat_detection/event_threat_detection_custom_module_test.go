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

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
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
		// Retry addCustomModule up to 5 times for each module
		var id string
		var err error
		for retry := 0; retry < 5; retry++ {
			id, err = addCustomModule()
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		if id != "" {
			sharedModules = append(sharedModules, id)
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

	// Define the custom module configuration
	customModule := &securitycentermanagementpb.EventThreatDetectionCustomModule{
		Config: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"metadata": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"description":    structpb.NewStringValue("Sample custom module description"),
						"severity":       structpb.NewStringValue("HIGH"),
						"recommendation": structpb.NewStringValue("Sample recommendation"),
					},
				}),
				"complianceStatus": structpb.NewStringValue("COMPLIANT"),
			},
		},
		DisplayName:     displayName,
		EnablementState: securitycentermanagementpb.EventThreatDetectionCustomModule_ENABLED,
		Type:            "CONFIGURABLE_BAD_IP",
	}

	req := &securitycentermanagementpb.CreateEventThreatDetectionCustomModuleRequest{
		Parent:                         parent,
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
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer
		var createModulePath = ""

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		// Call Create
		err := createEventThreatDetectionCustomModule(&buf, parent)

		if err != nil {
			r.Errorf("createCustomModule() had error: %v", err)
			return
		}

		got := buf.String()

		if got == "" {
			r.Errorf("createEventThreatDetectionCustomModule() returned an empty string")
			return
		}

		fmt.Printf("Response: %v\n", got)

		parts := strings.Split(got, ":")
		if len(parts) > 0 {
			createModulePath = parts[len(parts)-1]
		}

		AddModuleToCleanup(extractCustomModuleID(createModulePath))

		if !strings.Contains(got, orgID) {
			r.Errorf("createCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestGetCustomModule verifies the Get functionality
func TestGetEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		// Call Get
		err := getEventThreatDetectionCustomModule(&buf, parent, moduleID)

		if err != nil {
			r.Errorf("getEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()
		fmt.Printf("Response: %v\n", got)

		if !strings.Contains(got, orgID) {
			r.Errorf("getEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestUpdateCustomModule verifies the Update functionality
func TestUpdateEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
		// Call Update
		err := updateEventThreatDetectionCustomModule(&buf, parent, moduleID)

		if err != nil {
			r.Errorf("updateEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()

		if !strings.Contains(got, orgID) {
			r.Errorf("updateCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestDeleteCustomModule verifies the List functionality
func TestDeleteEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		// Create a dedicated module to delete
		id, err := addCustomModule()
		if err != nil {
			r.Errorf("addCustomModule() had error: %v", err)
			return
		}
		AddModuleToCleanup(id)

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		err = deleteEventThreatDetectionCustomModule(&buf, parent, id)

		if err != nil {
			r.Errorf("deleteEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()

		if !strings.Contains(got, id) {
			r.Errorf("deleteEventThreatDetectionCustomModule() got: %s want %s", got, id)
		}
	})
}

// TestListEtdCustomModule verifies the List functionality
func TestListEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		err := listEventThreatDetectionCustomModule(&buf, parent)

		if err != nil {
			r.Errorf("listEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()
		fmt.Printf("Response: %v\n", got)

		if !strings.Contains(got, orgID) {
			r.Errorf("listEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestListEffectiveEtdCustomModule verifies the List functionality
func TestListEffectiveEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		err := listEffectiveEventThreatDetectionCustomModule(&buf, parent)

		if err != nil {
			r.Errorf("listEffectiveEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()
		fmt.Printf("Response: %v\n", got)

		if !strings.Contains(got, orgID) {
			r.Errorf("listEffectiveEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestGetEffectiveEtdCustomModule verifies the Get functionality
func TestGetEffectiveEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		// Call Get
		err := getEffectiveEventThreatDetectionCustomModule(&buf, parent, moduleID)

		if err != nil {
			r.Errorf("getEffectiveEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()
		fmt.Printf("Response: %v\n", got)

		if !strings.Contains(got, orgID) {
			r.Errorf("getEffectiveEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestListDescendantEtdCustomModule verifies the List functionality
func TestListDescendantEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		err := listDescendantEventThreatDetectionCustomModule(&buf, parent)

		if err != nil {
			r.Errorf("listDescendantEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()
		fmt.Printf("Response: %v\n", got)

		if !strings.Contains(got, orgID) {
			r.Errorf("listDescendantEventThreatDetectionCustomModule() got: %s want %s", got, orgID)
		}
	})
}

// TestValidateEtdCustomModule verifies the List functionality
func TestValidateEtdCustomModule(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		var buf bytes.Buffer

		moduleID = getRandomSharedModule()
		if moduleID == "" {
			r.Errorf("No shared modules available")
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

		err := validateEventThreatDetectionCustomModule(&buf, parent)

		if err != nil {
			r.Errorf("validateEventThreatDetectionCustomModule() had error: %v", err)
			return
		}

		got := buf.String()
		fmt.Printf("Response: %v\n", got)

		// Check that the response indicates successful validation
		expectedMessage := "Validation successful: No errors found."

		if !strings.Contains(got, expectedMessage) {
			r.Errorf("validateEventThreatDetectionCustomModule() got: %s want %s", got, expectedMessage)
		}
	})
}
