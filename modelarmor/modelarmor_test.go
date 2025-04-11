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

func TestGetProjectFloorSettings(t *testing.T) {
	tc := testutil.SystemTest(t)

	var b bytes.Buffer
	if _, err := getProjectFloorSettings(&b, tc.ProjectID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved floor setting:"; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestGetOrganizationFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatal(err)
	}

	organizationID := os.Getenv("GOLANG_SAMPLES_ORGANIZATION_ID")
	var b bytes.Buffer
	if _, err := getOrganizationFloorSettings(&b, organizationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved org floor setting:"; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}

func TestGetFolderFloorSettings(t *testing.T) {
	// Load the test.env file
	err := godotenv.Load("./testdata/env/test.env")
	if err != nil {
		t.Fatal(err.Error())
	}

	folderID := os.Getenv("GOLANG_SAMPLES_FOLDER_ID")
	var b bytes.Buffer
	if _, err := getFolderFloorSettings(&b, folderID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved folder floor setting: "; !strings.Contains(got, want) {
		t.Errorf("getFloorSettings: expected %q to contain %q", got, want)
	}
}
