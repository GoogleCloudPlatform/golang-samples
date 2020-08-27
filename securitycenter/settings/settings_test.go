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
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func setup(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	return orgID
}

func TestEnableAssetDiscovery(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		orgID := setup(t)
		buf := new(bytes.Buffer)

		err := enableAssetDiscovery(buf, orgID)

		if err != nil {
			r.Errorf("enableAssetDiscovery(%s) had error: %v", orgID, err)
			return
		}

		got := buf.String()
		if want := "Asset discovery on? true"; !strings.Contains(got, want) {
			r.Errorf("enableAssetDiscovery(%s) got: %s want %s", orgID, got, want)
		}
		if !strings.Contains(got, orgID) {
			r.Errorf("enableAssetDiscovery(%s) got: %s want %s", orgID, got, orgID)
		}
	})
}

func TestGetOrgSettings(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		orgID := setup(t)
		buf := new(bytes.Buffer)

		err := getOrgSettings(buf, orgID)

		if err != nil {
			r.Errorf("getOrgSettings(%s) had error: %v", orgID, err)
			return
		}

		got := buf.String()
		if want := "Asset Discovery on? "; !strings.Contains(got, want) {
			r.Errorf("getOrgSettings(%s) got: %s want %s", orgID, got, want)
		}
		if !strings.Contains(got, orgID) {
			r.Errorf("getOrgSettings(%s) got: %s want %s", orgID, got, orgID)
		}
	})
}
