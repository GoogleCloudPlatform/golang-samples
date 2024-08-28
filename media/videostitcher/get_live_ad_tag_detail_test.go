// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package videostitcher

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetLiveAdTagDetail(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	slateID := fmt.Sprintf("%s-%s", slateIDPrefix, uuid)
	slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
	liveConfigID := fmt.Sprintf("%s-%s", liveConfigIDPrefix, uuid)
	liveConfigName := fmt.Sprintf("projects/%s/locations/%s/liveConfigs/%s", tc.ProjectID, location, liveConfigID)

	createTestSlate(slateID, t)
	createTestLiveConfig(slateID, liveConfigID, t)
	t.Cleanup(func() {
		// Can't delete live sessions
		deleteTestLiveConfig(liveConfigName, t)
		deleteTestSlate(slateName, t)
	})

	sessionID, playURI := createTestLiveSession(liveConfigID, t)

	// No list or delete methods for live sessions

	// Ad tag details

	// To get ad tag details, you need to curl the main manifest and
	// a rendition first. This supplies media player information to the API.
	//
	// Get the playURI first. The last line of the response will contain a
	// renditions location. Curl the live session name with the rendition
	// location appended.

	renditions, err := getPlayURI(playURI)
	if err != nil {
		t.Fatalf("getPlayURI err: %v", err)
	}

	// playURI will be in the following format:
	// https://videostitcher.googleapis.com/v1/projects/{project}/locations/{location}/liveSessions/{session-id}/manifest.m3u8?signature=...
	// Replace manifest.m3u8?signature=... with the renditions location.

	err = curlRendition(playURI, renditions[0])
	if err != nil {
		t.Fatalf("curlRendition err: %v", err)
	}

	// List the ad tag details for a given live session. This is the only way to get an
	// ad tag detail.
	adTagDetailsNamePrefix := fmt.Sprintf("/locations/%s/liveSessions/%s/liveAdTagDetails/", location, sessionID)
	if err := listLiveAdTagDetails(&buf, tc.ProjectID, sessionID); err != nil {
		t.Fatalf("listLiveAdTagDetails got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, adTagDetailsNamePrefix) {
		t.Fatalf("listLiveAdTagDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsNamePrefix)
	}
	strSlice := strings.Split(got, "/")
	adTagDetailsID := strSlice[len(strSlice)-1]
	adTagDetailsID = strings.TrimRight(adTagDetailsID, "\n")
	adTagDetailsName := fmt.Sprintf("/locations/%s/liveSessions/%s/liveAdTagDetails/%s", location, sessionID, adTagDetailsID)
	if !strings.Contains(got, adTagDetailsName) {
		t.Errorf("listLiveAdTagDetails got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsName)
	}
	buf.Reset()

	// Get the specified ad tag detail for a given live session.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := getLiveAdTagDetail(&buf, tc.ProjectID, sessionID, adTagDetailsID); err != nil {
			r.Errorf("getLiveAdTagDetail got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, adTagDetailsName) {
			r.Errorf("getLiveAdTagDetail got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, adTagDetailsName)
		}
	})
	buf.Reset()
}
