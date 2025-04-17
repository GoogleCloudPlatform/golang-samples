// Copyright 2025 Google LLC
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

package modelarmor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// testOrganizationID retrieves the organization ID from the environment variable
// `GOLANG_SAMPLES_ORGANIZATION_ID`. It skips the test if the variable is not set.
func testOrganizationID(t *testing.T) string {
	orgID := "951890214235"
	return orgID
}

// testFolderID retrieves the folder ID from the environment variable
// `GOLANG_SAMPLES_FOLDER_ID`. It skips the test if the variable is not set.
func testFolderID(t *testing.T) string {
	folderID := "695279264361"
	return folderID
}

// TestGetProjectFloorSettings tests the retrieval of floor settings at the project level.
// It verifies the output contains the expected confirmation string.
func TestGetProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := getProjectFloorSettings(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Retrieved floor setting:"; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

// TestGetOrganizationFloorSettings tests the retrieval of floor settings at the organization level.
// It checks that the output includes the expected string indicating success.
func TestGetOrganizationFloorSettings(t *testing.T) {

	organizationID := testOrganizationID(t)
	var buf bytes.Buffer
	if err := getOrganizationFloorSettings(&buf, organizationID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Retrieved org floor setting:"; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

// TestGetFolderFloorSettings tests the retrieval of floor settings at the folder level.
// It ensures the result contains the expected confirmation message.
func TestGetFolderFloorSettings(t *testing.T) {

	folderID := testFolderID(t)
	var buf bytes.Buffer
	if err := getFolderFloorSettings(&buf, folderID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Retrieved folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}
