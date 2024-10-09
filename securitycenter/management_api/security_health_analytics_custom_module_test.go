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
	"strings"
	"testing"
	// "time"

	securitycenter "cloud.google.com/go/securitycentermanagement/apiv1"
	securitycenterpb "cloud.google.com/go/securitycentermanagement/apiv1/securitycentermanagementpb"
	"github.com/google/uuid"
	// "github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	// exprpb "google.golang.org/genproto/googleapis/type/expr"
)

var orgID = ""

// setup initializes orgID for testing
func setup(t *testing.T) string {
	orgID = "1081635000895" // Replace with your actual organization ID
	return orgID
}

// func addCustomModule(t *testing.T, customModuleID string) error {
// 	orgID := setup(t)
	
// 	fmt.Printf("Custom Module ID: %v\n", customModuleID)

// 	ctx := context.Background()	
// 	client, err := securitycenter.NewClient(ctx)
	
// 	fmt.Println("Error from client:", err)

// 	if err != nil {
// 		return fmt.Errorf("securitycenter.NewClient: %w", err)
// 	}
// 	defer client.Close()

// 	customModule := &securitycenterpb.SecurityHealthAnalyticsCustomModule{
// 		// DisplayName: "CustomModule for testing",
// 		// EnablementState: securitycenterpb.SecurityHealthAnalyticsCustomModule_ENABLED,
// 		DisplayName: "CustomModule for testing",
//         EnablementState: securitycenterpb.SecurityHealthAnalyticsCustomModule_ENABLED,
//         CustomConfig: &securitycenterpb.CustomConfig{
//             Predicate: &exprpb.Expr{
//                 Expression: "resource.type == \"gce_instance\"",
//             },
//             ResourceSelector: &securitycenterpb.CustomConfig_ResourceSelector{
//                 ResourceTypes: []string{"gce_instance"},
//             },
//             Severity:     securitycenterpb.CustomConfig_CRITICAL,
//             Description:  "Detect high severity issues on GCE instances.",
//             Recommendation: "Ensure proper configurations for GCE instances.",
//         },
// 	}

// 	req := &securitycenterpb.CreateSecurityHealthAnalyticsCustomModuleRequest{
// 		Parent:                      fmt.Sprintf("organizations/%s/locations/global", orgID),
// 		SecurityHealthAnalyticsCustomModule: customModule,
// 	}

// 	_, err0 := client.CreateSecurityHealthAnalyticsCustomModule(ctx, req)
// 	if err0 != nil {
// 		return fmt.Errorf("Failed to create CustomModule: %w", err0)
// 	}
// 	return nil
// }

func cleanupCustomModule(t *testing.T, customModuleID string) error {
	orgID := setup(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.DeleteSecurityHealthAnalyticsCustomModuleRequest{
		Name: fmt.Sprintf("organizations/%s/locations/global/securityHealthAnalyticsCustomModules/%s", orgID, customModuleID),
	}

	if err := client.DeleteSecurityHealthAnalyticsCustomModule(ctx, req); err != nil {
		return fmt.Errorf("failed to delete CustomModule: %w", err)
	}

	return nil
}

// TestCreateCustomModule verifies the Create functionality
 func TestCreateCustomModule(t *testing.T) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	rand, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Issue generating id.")
		return
	}
	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)	
	customModuleID := "custommodule_id_" + rand.String()

	// Call Create
	err = CreateSecurityHealthAnalyticsCustomModule(buf, parent, customModuleID)

	if err != nil {
		t.Fatalf("createCustomModule() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, customModuleID) {
		t.Fatalf("createCustomModule() got: %s want %s", got, customModuleID)
	}

	// Cleanup
	// cleanupCustomModule(t, customModuleID)
}

// TestGetCustomModule verifies the Get functionality
func TestGetCustomModule(t *testing.T) {
	orgID := setup(t)

	buf := new(bytes.Buffer)

	// Create Test CustomModule
	_, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Issue generating id.")
		return
	}
	// cusModID := "random-custommodule-id-" + rand.String()
	customModuleID := "10829723016802655264"

	// if err := addCustomModule(t, customModuleID); err != nil {
	// 	t.Fatalf("Could not setup test environment: %v", err)
	// 	return
	// }

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	// Call Get
	err = getSecurityHealthAnalyticsCustomModule(buf, parent, customModuleID)

	if err != nil {
		t.Fatalf("getSecurityHealthAnalyticsCustomModule() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, customModuleID) {
		t.Fatalf("getSecurityHealthAnalyticsCustomModule() got: %s want %s", got, customModuleID)
	}

	// Cleanup
	// cleanupCustomModule(t, customModuleID)
}
