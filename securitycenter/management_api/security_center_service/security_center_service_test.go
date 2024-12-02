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

package security_center_service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

var orgID = ""

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

// TestGetSecurityCenterService verifies the Get functionality
func TestGetSecurityCenterService(t *testing.T) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	service := "event-threat-detection"
	// Call the function
	err := getSecurityCenterService(&buf, parent, service)

	if err != nil {
		t.Fatalf("getSecurityCenterService() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	// Check if the output contains the service name
	if !strings.Contains(got, service) {
		t.Fatalf("getSecurityCenterService() got: %s want %s", got, service)
	}
}

// TestListSecurityCenterService verifies the List functionality
func TestListSecurityCenterService(t *testing.T) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)

	err := listSecurityCenterService(&buf, parent)

	if err != nil {
		t.Fatalf("listSecurityCenterService() had error: %v", err)
		return
	}

	got := buf.String()
	fmt.Printf("Response: %v\n", got)

	if !strings.Contains(got, orgID) {
		t.Fatalf("listSecurityCenterService() got: %s want %s", got, orgID)
	}
}

// TestUpdateSecurityCenterService verifies the Update functionality
func TestUpdateSecurityCenterService(t *testing.T) {
	var buf bytes.Buffer

	parent := fmt.Sprintf("organizations/%s/locations/global", orgID)
	service := "event-threat-detection"
	// Call the function
	err := updateSecurityCenterService(&buf, parent, service)

	if err != nil {
		t.Fatalf("updateSecurityCenterService() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, orgID) {
		t.Fatalf("updateSecurityCenterService() got: %s want %s", got, orgID)
	}
}
