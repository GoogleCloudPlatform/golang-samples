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
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// testOrganizationID retrieves the organization ID from the environment variable
// `GOLANG_SAMPLES_ORGANIZATION_ID`. It skips the test if the variable is not set.
func testOrganizationID(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_ORGANIZATION_ID")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_SAMPLES_ORGANIZATION_ID")
	}

	return v
}

// testFolderID retrieves the folder ID from the environment variable
// `GOLANG_SAMPLES_FOLDER_ID`. It skips the test if the variable is not set.
func testFolderID(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_FOLDER_ID")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_SAMPLES_FOLDER_ID")
	}

	return v
}

// testLocation returns the location for testing from the environment variable.
// Skips the test if the environment variable is not set.
func testLocation(t *testing.T) string {
	t.Helper()

	v := os.Getenv("GOLANG_SAMPLES_LOCATION")
	if v == "" {
		t.Skip("testIamUser: missing GOLANG_SAMPLES_LOCATION")
	}

	return v
}

// TestUpdateFolderFloorSettings tests updating floor settings for a specific folder.
// It verifies that the output buffer contains a confirmation message indicating a successful update.
func TestUpdateFolderFloorSettings(t *testing.T) {
	folderID := testFolderID(t)
	var b bytes.Buffer
	if err := updateFolderFloorSettings(&b, folderID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateFolderFloorSettings: expected %q to contain %q", got, want)
	}
}

// TestUpdateOrganizationFloorSettings tests updating floor settings for a specific organization.
// It ensures the output buffer includes a success message confirming the update.
func TestUpdateOrganizationFloorSettings(t *testing.T) {
	organizationID := testOrganizationID(t)
	locationID := testLocation(t)
	var buf bytes.Buffer
	if err := updateOrganizationFloorSettings(&buf, organizationID, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Updated org floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateOrganizationFloorSettings: expected %q to contain %q", got, want)
	}
}

// TestUpdateProjectFloorSettings tests updating floor settings for a specific project.
// It checks that the resulting output includes the expected confirmation message.
func TestUpdateProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	var buf bytes.Buffer
	if err := updateProjectFloorSettings(&buf, tc.ProjectID, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := buf.String(), "Updated project floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateProjectFloorSettings: expected %q to contain %q", got, want)
	}
}
