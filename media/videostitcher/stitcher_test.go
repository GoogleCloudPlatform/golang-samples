// Copyright 2022 Google LLC
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

package videostitcher

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	deleteCdnKeyResponse = "Deleted CDN key"
	deleteSlateResponse  = "Deleted slate"
	location             = "us-central1" // All samples use this location
	slateID              = "my-go-test-slate"
)

var bucketName string
var slateURI string
var updatedSlateURI string
var projectNumber string

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following API on the test project:
// *   Video Stitcher API
// *   Cloud Resource Manager API (needed for project number translation)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName = "cloud-samples-data/media/"
	slateURI = "https://storage.googleapis.com/" + bucketName + "ForBiggerEscapes.mp4"
	updatedSlateURI = "https://storage.googleapis.com/" + bucketName + "ForBiggerJoyrides.mp4"

	// Get the project number
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		t.Fatalf("cloudresourcemanager.NewService: %v", err)
	}
	project, err := cloudresourcemanagerClient.Projects.Get(tc.ProjectID).Do()
	if err != nil {
		t.Fatalf("cloudresourcemanagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber = strconv.FormatInt(project.ProjectNumber, 10)
}

// testSlates tests major operations on slates. Create, list, update,
// and get operations check if the slate resource name is returned. The
// delete operation checks for a hard-coded string response.
func TestSlates(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Test setup

	// Delete the default slate if it exists
	deleteSlate(buf, tc.ProjectID, slateID)
	defer deleteSlate(buf, tc.ProjectID, slateID)

	// Tests

	// Create a new slate.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", projectNumber, location, slateID)
		if err := createSlate(buf, tc.ProjectID, slateID, slateURI); err != nil {
			r.Errorf("createSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			r.Errorf("createSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
		}
	})
	buf.Reset()

	// List the slates for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
		if err := listSlates(buf, tc.ProjectID); err != nil {
			r.Errorf("listSlates got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			r.Errorf("listSlates got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
		}
	})
	buf.Reset()

	// Update an existing slate.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
		if err := updateSlate(buf, tc.ProjectID, slateID, updatedSlateURI); err != nil {
			r.Errorf("updateSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			r.Errorf("updateSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
		}
		if got := buf.String(); !strings.Contains(got, updatedSlateURI) {
			r.Errorf("updateSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, updatedSlateURI)
		}
	})
	buf.Reset()

	// Get the updated slate.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
		if err := getSlate(buf, tc.ProjectID, slateID); err != nil {
			r.Errorf("getSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			r.Errorf("getSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
		}
	})

	// Delete the slate.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteSlate(buf, tc.ProjectID, slateID); err != nil {
			r.Errorf("deleteSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteSlateResponse) {
			r.Errorf("deleteSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteSlateResponse)
		}
	})
}
