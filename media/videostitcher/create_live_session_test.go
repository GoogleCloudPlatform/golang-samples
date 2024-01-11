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

func TestCreateLiveSession(t *testing.T) {
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
		deleteTestLiveConfig(liveConfigName, t)
		deleteTestSlate(slateName, t)
	})

	// Create a new live session and return the play URI.
	sessionPrefix := fmt.Sprintf("locations/%s/liveSessions/", location)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := createLiveSession(&buf, tc.ProjectID, liveConfigID); err != nil {
			r.Errorf("createLiveSession got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, sessionPrefix) {
			r.Errorf("createLiveSession got: %v Want to contain: %v", got, sessionPrefix)
		}
	})
}
