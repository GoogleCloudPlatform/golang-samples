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
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	location            = "us-central1" // All samples use this location
	slateID             = "my-go-test-slate"
	deleteSlateResponse = "Deleted slate"

	deleteCDNKeyResponse = "Deleted CDN key"
	mediaCDNKeyID        = "my-go-test-media-cdn"
	cloudCDNKeyID        = "my-go-test-cloud-cdn"
	akamaiCDNKeyID       = "my-go-test-akamai-cdn"
	hostname             = "cdn.example.com"
	updatedHostname      = "updated.example.com"
	keyName              = "my-key"
)

var bucketName string
var slateURI string
var updatedSlateURI string
var projectNumber string
var vodURI string
var liveURI string

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
	vodURI = "https://storage.googleapis.com/" + bucketName + "hls-vod/manifest.m3u8"
	liveURI = "https://storage.googleapis.com/" + bucketName + "hls-live/manifest.m3u8"

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
	t.Skip("see GoogleCloudPlatform/golang-samples#2919")
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
			t.Fatalf("createSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			t.Fatalf("createSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
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

// TestCDNKeys tests major operations on CDN keys. Create, list, update,
// and get operations check if the CDN key resource name is returned. The
// delete operation checks for a hard-coded string response.
func TestCDNKeys(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#2920")
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Test setup

	// Delete the Media CDN key if it exists.
	if err := getCDNKey(buf, tc.ProjectID, mediaCDNKeyID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteCDNKey(buf, tc.ProjectID, mediaCDNKeyID); err != nil {
				r.Errorf("deleteCDNKey got err: %v", err)
			}
		})
	}

	// Delete the Cloud CDN key if it exists.
	if err := getCDNKey(buf, tc.ProjectID, cloudCDNKeyID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteCDNKey(buf, tc.ProjectID, cloudCDNKeyID); err != nil {
				r.Errorf("deleteCDNKey got err: %v", err)
			}
		})
	}

	// Delete the Akamai CDN key if it exists.
	if err := getCDNKey(buf, tc.ProjectID, akamaiCDNKeyID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteCDNKey(buf, tc.ProjectID, akamaiCDNKeyID); err != nil {
				r.Errorf("deleteCDNKey got err: %v", err)
			}
		})
	}

	// Tests
	// Media CDN tests
	mediaCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}

	// Create a new Media CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectNumber, location, mediaCDNKeyID)
		if err := createCDNKey(buf, tc.ProjectID, mediaCDNKeyID, hostname, keyName, mediaCDNPrivateKey, true); err != nil {
			t.Fatalf("createCDNKey (Media CDN) got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			t.Fatalf("createCDNKey (Media CDN) got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// List the CDN keys for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, mediaCDNKeyID)
		if err := listCDNKeys(buf, tc.ProjectID); err != nil {
			r.Errorf("listCDNKeys got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("listCDNKeys got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// Update an existing CDN key.
	updatedMediaCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, mediaCDNKeyID)
		if err := updateCDNKey(buf, tc.ProjectID, mediaCDNKeyID, updatedHostname, keyName, updatedMediaCDNPrivateKey, true); err != nil {
			r.Errorf("updateCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("updateCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// Get the updated CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, mediaCDNKeyID)
		if err := getCDNKey(buf, tc.ProjectID, mediaCDNKeyID); err != nil {
			r.Errorf("getCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("getCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})

	// Delete the CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteCDNKey(buf, tc.ProjectID, mediaCDNKeyID); err != nil {
			r.Errorf("deleteCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteCDNKeyResponse) {
			r.Errorf("deleteCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteCDNKeyResponse)
		}
	})

	// Cloud CDN tests

	// Create a new Cloud CDN key.
	cloudCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectNumber, location, cloudCDNKeyID)
		if err := createCDNKey(buf, tc.ProjectID, cloudCDNKeyID, hostname, keyName, cloudCDNPrivateKey, false); err != nil {
			t.Fatalf("createCDNKey (Cloud CDN) got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			t.Fatalf("createCDNKey (Cloud CDN) got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// List the CDN keys for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, cloudCDNKeyID)
		if err := listCDNKeys(buf, tc.ProjectID); err != nil {
			r.Errorf("listCDNKeys got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("listCDNKeys got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// Update an existing CDN key.
	updatedCloudCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, cloudCDNKeyID)
		if err := updateCDNKey(buf, tc.ProjectID, cloudCDNKeyID, updatedHostname, keyName, updatedCloudCDNPrivateKey, false); err != nil {
			r.Errorf("updateCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("updateCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// Get the updated CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, cloudCDNKeyID)
		if err := getCDNKey(buf, tc.ProjectID, cloudCDNKeyID); err != nil {
			r.Errorf("getCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("getCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})

	// Delete the CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteCDNKey(buf, tc.ProjectID, cloudCDNKeyID); err != nil {
			r.Errorf("deleteCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteCDNKeyResponse) {
			r.Errorf("deleteCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteCDNKeyResponse)
		}
	})

	// Akamai tests

	// Create a new Akamai CDN key.
	akamaiTokenKey, err := getUUID64()
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectNumber, location, akamaiCDNKeyID)
		if err := createCDNKeyAkamai(buf, tc.ProjectID, akamaiCDNKeyID, hostname, akamaiTokenKey); err != nil {
			t.Fatalf("createCDNKeyAkamai got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			t.Fatalf("createCDNKeyAkamai got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// List the CDN keys for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, akamaiCDNKeyID)
		if err := listCDNKeys(buf, tc.ProjectID); err != nil {
			r.Errorf("listCDNKeys got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("listCDNKeys got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// Update an existing CDN key.
	updatedAkamaiTokenKey, err := getUUID64()
	if err != nil {
		t.Fatalf("uuid err: %v", err)
	}

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, akamaiCDNKeyID)
		if err := updateCDNKeyAkamai(buf, tc.ProjectID, akamaiCDNKeyID, updatedHostname, updatedAkamaiTokenKey); err != nil {
			r.Errorf("updateCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("updateCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})
	buf.Reset()

	// Get the updated CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		cdnKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, akamaiCDNKeyID)
		if err := getCDNKey(buf, tc.ProjectID, akamaiCDNKeyID); err != nil {
			r.Errorf("getCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cdnKeyName) {
			r.Errorf("getCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, cdnKeyName)
		}
	})

	// Delete the CDN key.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteCDNKey(buf, tc.ProjectID, akamaiCDNKeyID); err != nil {
			r.Errorf("deleteCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteCDNKeyResponse) {
			r.Errorf("deleteCDNKey got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteCDNKeyResponse)
		}
	})
}

// testVodSessions tests major operations on VOD sessions. Create and get
// operations check if the session name is returned. List and delete methods
// are not supported for VOD sessions.
// The test lists and gets ad tag and stitch details as well.
func TestVodSessions(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#2917")
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	sessionID := ""

	// Create a new VOD session.
	sessionPrefix := fmt.Sprintf("projects/%s/locations/%s/vodSessions/", projectNumber, location)
	if err := createVodSession(buf, tc.ProjectID, vodURI); err != nil {
		t.Fatalf("createVodSession got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, sessionPrefix) {
		t.Fatalf("createVodSession got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, sessionPrefix)
	}
	strSlice := strings.Split(got, "/")
	sessionID = strSlice[len(strSlice)-1]
	buf.Reset()

	// Get the VOD session.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		sessionName := fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s", projectNumber, location, sessionID)
		if err := getVodSession(buf, tc.ProjectID, sessionID); err != nil {
			r.Errorf("getVodSession got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, sessionName) {
			r.Errorf("getVodSession got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, sessionName)
		}
	})
	buf.Reset()

	// No list or delete methods for VOD sessions

	// Ad tag details

	// List the ad tag details for a given VOD session.
	adTagDetailsNamePrefix := fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s/vodAdTagDetails/", projectNumber, location, sessionID)
	if err := listVodAdTagDetails(buf, tc.ProjectID, sessionID); err != nil {
		t.Errorf("listVodAdTagDetails got err: %v", err)
	}
	got = buf.String()

	if !strings.Contains(got, adTagDetailsNamePrefix) {
		t.Errorf("listVodAdTagDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsNamePrefix)
	}
	strSlice = strings.Split(got, "/")
	adTagDetailsID := strSlice[len(strSlice)-1]
	adTagDetailsID = strings.TrimRight(adTagDetailsID, "\n")
	adTagDetailsName := fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s/vodAdTagDetails/%s", projectNumber, location, sessionID, adTagDetailsID)
	if !strings.Contains(got, adTagDetailsName) {
		t.Errorf("listVodAdTagDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsName)
	}
	buf.Reset()

	// Get the specified ad tag detail for a given VOD session.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := getVodAdTagDetail(buf, tc.ProjectID, sessionID, adTagDetailsID); err != nil {
			r.Errorf("getVodAdTagDetail got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, adTagDetailsName) {
			r.Errorf("getVodAdTagDetail got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsName)
		}
	})
	buf.Reset()

	// Stitch details

	// List the stitch details for a given VOD session.
	stitchDetailsNamePrefix := fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s/vodStitchDetails/", projectNumber, location, sessionID)
	if err := listVodStitchDetails(buf, tc.ProjectID, sessionID); err != nil {
		t.Errorf("listVodStitchDetails got err: %v", err)
	}
	got = buf.String()

	if !strings.Contains(got, stitchDetailsNamePrefix) {
		t.Errorf("listVodStitchDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, stitchDetailsNamePrefix)
	}
	strSlice = strings.Split(got, "/")
	stitchDetailsID := strSlice[len(strSlice)-1]
	stitchDetailsID = strings.TrimRight(stitchDetailsID, "\n")
	stitchDetailsName := fmt.Sprintf("projects/%s/locations/%s/vodSessions/%s/vodStitchDetails/%s", projectNumber, location, sessionID, stitchDetailsID)
	if !strings.Contains(got, stitchDetailsName) {
		t.Errorf("listVodStitchDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, stitchDetailsName)
	}
	buf.Reset()

	// Get the specified VOD stitch detail for a given VOD session.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := getVodStitchDetail(buf, tc.ProjectID, sessionID, stitchDetailsID); err != nil {
			r.Errorf("getVodStitchDetail got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, stitchDetailsName) {
			r.Errorf("getVodStitchDetail got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, stitchDetailsName)
		}
	})
	buf.Reset()
}

// testLiveSessions tests major operations on live sessions. Create and get
// operations check if the session name is returned. List and delete methods
// are not supported for live sessions. The test lists and gets ad tag details.
func TestLiveSessions(t *testing.T) {
	t.Skip("see GoogleCloudPlatform/golang-samples#2918")
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	sessionID := ""

	// Test setup

	deleteSlate(buf, tc.ProjectID, slateID)
	defer deleteSlate(buf, tc.ProjectID, slateID)

	// Tests

	// Create a new slate.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", projectNumber, location, slateID)
		if err := createSlate(buf, tc.ProjectID, slateID, slateURI); err != nil {
			t.Fatalf("createSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			t.Fatalf("createSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
		}
	})
	buf.Reset()

	// Create a new live session and return the play URI.
	sessionPrefix := fmt.Sprintf("projects/%s/locations/%s/liveSessions/", projectNumber, location)
	if err := createLiveSession(buf, tc.ProjectID, liveURI, slateID); err != nil {
		t.Fatalf("createLiveSession got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, sessionPrefix) {
		t.Fatalf("createLiveSession got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, sessionPrefix)
	}
	buf.Reset()

	re := regexp.MustCompile(fmt.Sprintf("Live session:.%s(.*)", sessionPrefix))
	match := re.FindAllStringSubmatch(string(got), -1)
	if len(match) == 1 && len(match[0]) == 2 {
		sessionID = match[0][1]
	} else {
		t.Errorf("\nSession ID not found in %s\n", got)
	}

	re = regexp.MustCompile(`Play URI:.(.*)`)
	match = re.FindAllStringSubmatch(string(got), -1)
	playURI := ""
	if len(match) == 1 && len(match[0]) == 2 {
		playURI = match[0][1]
	} else {
		t.Errorf("\nPlay URI not found %s\n", string(got))
	}

	// Get the live session.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		sessionName := fmt.Sprintf("projects/%s/locations/%s/liveSessions/%s", projectNumber, location, sessionID)
		if err := getLiveSession(buf, tc.ProjectID, sessionID); err != nil {
			r.Errorf("getLiveSession got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, sessionName) {
			r.Errorf("getLiveSession got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, sessionName)
		}
	})
	buf.Reset()

	// No list or delete methods for live sessions

	// Ad tag details

	// To get ad tag details, you need to curl the main manifest and
	// a rendition first. This supplies media player information to the API.
	//
	// Curl the playURI first. The last line of the response will contain a
	// renditions location. Curl the live session name with the rendition
	// location appended.

	resp, err := http.Get(playURI)
	if err != nil {
		t.Errorf("\nError getting the play URI\n")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll: %v", err)
	}

	re = regexp.MustCompile(`renditions/.*`)
	renditions := re.FindStringSubmatch(string(body))
	if len(renditions) == 0 {
		t.Errorf("\nRenditions not found in body: %s\n", string(body))
	}

	// playURI will be in the following format:
	// https://videostitcher.googleapis.com/v1/projects/{project}/locations/{location}/liveSessions/{session-id}/manifest.m3u8?signature=...
	// Replace manifest.m3u8?signature=... with the renditions location.

	re = regexp.MustCompile(`manifest.m3u8.*`)
	renditionURI := re.ReplaceAllString(playURI, renditions[0])
	resp, err = http.Get(renditionURI)
	if err != nil {
		t.Errorf("\nError getting the rendition URI\n")
	}
	defer resp.Body.Close()

	// List the ad tag details for a given live session.
	adTagDetailsNamePrefix := fmt.Sprintf("projects/%s/locations/%s/liveSessions/%s/liveAdTagDetails/", projectNumber, location, sessionID)
	if err := listLiveAdTagDetails(buf, tc.ProjectID, sessionID); err != nil {
		t.Errorf("listLiveAdTagDetails got err: %v", err)
	}
	got = buf.String()

	if !strings.Contains(got, adTagDetailsNamePrefix) {
		t.Errorf("listLiveAdTagDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsNamePrefix)
	}
	strSlice := strings.Split(got, "/")
	adTagDetailsID := strSlice[len(strSlice)-1]
	adTagDetailsID = strings.TrimRight(adTagDetailsID, "\n")
	adTagDetailsName := fmt.Sprintf("projects/%s/locations/%s/liveSessions/%s/liveAdTagDetails/%s", projectNumber, location, sessionID, adTagDetailsID)
	if !strings.Contains(got, adTagDetailsName) {
		t.Errorf("listLiveAdTagDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsName)
	}
	buf.Reset()

	// Get the specified ad tag detail for a given live session.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := getLiveAdTagDetail(buf, projectNumber, sessionID, adTagDetailsID); err != nil {
			r.Errorf("getLiveAdTagDetail got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, adTagDetailsName) {
			r.Errorf("getLiveAdTagDetail got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsName)
		}
	})
	buf.Reset()
}

func getUUID64() (string, error) {
	u1, err1 := uuid.NewRandom()
	u2, err2 := uuid.NewRandom()
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("uuid err: %v, %v", err1, err2)
	}
	uuid := fmt.Sprintf("%s%s", u1.String(), u2.String())
	return strings.ReplaceAll(uuid, "-", ""), nil
}
