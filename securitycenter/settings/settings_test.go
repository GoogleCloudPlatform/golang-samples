// Copyright 2019 Google LLC
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

package settings

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func setup(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	return orgID
}

func TestEnableAssetDiscovery(t *testing.T) {
	orgID := setup(t)
	buf := new(bytes.Buffer)

	err := enableAssetDiscovery(buf, orgID)

	if err != nil {
		t.Fatalf("enableAssetDiscovery(%s) had error: %v", orgID, err)
	}

	got := buf.String()
	if want := "Asset discovery on? true"; !strings.Contains(got, want) {
		t.Errorf("enableAssetDiscovery(%s) got: %s want %s", orgID, got, want)
	}
	if !strings.Contains(got, orgID) {
		t.Errorf("enableAssetDiscovery(%s) got: %s want %s", orgID, got, orgID)
	}

}

func TestGetOrgSettings(t *testing.T) {
	orgID := setup(t)
	buf := new(bytes.Buffer)

	err := getOrgSettings(buf, orgID)

	if err != nil {
		t.Fatalf("getOrgSettings(%s) had error: %v", orgID, err)
	}

	got := buf.String()
	if want := "Asset Discovery on? "; !strings.Contains(got, want) {
		t.Errorf("getOrgSettings(%s) got: %s want %s", orgID, got, want)
	}
	if !strings.Contains(got, orgID) {
		t.Errorf("getOrgSettings(%s) got: %s want %s", orgID, got, orgID)
	}

}
