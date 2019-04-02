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

package assets

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func setup(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	return orgID
}

func TestListAllAssets(t *testing.T) {
	buf := new(bytes.Buffer)
	orgID := setup(t)
	got, err := listAllAssets(buf, orgID)
	if err != nil {
		t.Fatalf("listAllAssets(%s) failed: %v", orgID, err)
	}
	wanted := 59
	if got < 59 {
		t.Errorf("listAllAssets(%s) Not enough results: %d Wanted: %d", orgID, got, wanted)
	}
}

func TestListAllProjectAssets(t *testing.T) {
	buf := new(bytes.Buffer)
	orgID := setup(t)
	got, err := listAllProjectAssets(buf, orgID)
	if err != nil {
		t.Fatalf("listAllAssets(%s) failed: %v", orgID, err)
	}
	wanted := 3
	if got != 3 {
		t.Errorf("listAllAssets(%s) Unexpected number of results: %d Wanted: %d", orgID, got, wanted)
	}
}

func TestListAllProjectAssetsAtTime(t *testing.T) {
	orgID := setup(t)
	buf := new(bytes.Buffer)
	var nothingInstant = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	got, err := listAllProjectAssetsAtTime(buf, orgID, nothingInstant)
	if err != nil {
		t.Fatalf("listAllProjectAssetsAtTime(%s, %v) failed: %v", orgID, nothingInstant, err)
	}
	if got != 0 {
		t.Errorf("listAllProjectAssetsAtTime(%s, %v) Results not 0: %d", orgID, nothingInstant, got)
	}

	var somethingInstant = time.Date(2019, 3, 15, 0, 0, 0, 0, time.UTC)
	got, err = listAllProjectAssetsAtTime(buf, orgID, somethingInstant)
	if err != nil {
		t.Fatalf("listAllProjectAssetsAtTime(%s, %v) failed: %v", orgID, somethingInstant, err)
	}
	wanted := 3
	if got != 3 {
		t.Errorf("listAllProjectAssetsAtTime(%s, %v) Unexpected number of projects: %d Wanted: %d", orgID, somethingInstant, got, wanted)
	}
}

func TestListAllProjectAssetsAndStateChanges(t *testing.T) {
	buf := new(bytes.Buffer)
	orgID := setup(t)
	got, err := listAllProjectAssetsAndStateChanges(buf, orgID)
	if err != nil {
		t.Fatalf("listAllProjectAssetsAndStateChanges(%s) failed: %v", orgID, err)
	}
	wanted := 3
	if got != 3 {
		t.Errorf("listAllProjectAssetsAndStateChanges(%s) Unexpected number of results: %d Wanted: %d", orgID, got, wanted)
	}
}
