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
	"github.com/joho/godotenv"
)

func TestUpdateFolderFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatal(err.Error())
	}
	folderID := os.Getenv("GOLANG_SAMPLES_FOLDER_ID")
	var b bytes.Buffer
	if _, err := updateFolderFloorSettings(&b, folderID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateFolderFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestUpdateOrganizationFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatal(err.Error())
	}

	organizationID := os.Getenv("GOLANG_SAMPLES_ORGANIZATION_ID")
	var b bytes.Buffer
	if _, err := updateOrganizationFloorSettings(&b, organizationID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated org floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateOrganizationFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestUpdateProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)
	var b bytes.Buffer
	if _, err := updateProjectFloorSettings(&b, tc.ProjectID, "us-central1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated project floor setting: "; !strings.Contains(got, want) {
		t.Errorf("updateProjectFloorSettings: expected %q to contain %q", got, want)
	}
}
