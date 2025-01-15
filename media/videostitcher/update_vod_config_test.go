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

func TestUpdateVodConfig(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	vodConfigID := fmt.Sprintf("%s-%s", vodConfigIDPrefix, uuid)
	createTestVodConfig(vodConfigID, t)

	vodConfigName := fmt.Sprintf("projects/%s/locations/%s/vodConfigs/%s", tc.ProjectID, location, vodConfigID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := updateVodConfig(&buf, tc.ProjectID, vodConfigID, updatedVodURI); err != nil {
			r.Errorf("updateVodConfig got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, vodConfigName) {
			r.Errorf("updateVodConfig got: %v Want to contain: %v", got, vodConfigName)
		}
		if got := buf.String(); !strings.Contains(got, updatedVodURI) {
			r.Errorf("updateVodConfig got: %v Want to contain: %v", got, updatedVodURI)
		}
	})

	t.Cleanup(func() {
		deleteTestVodConfig(vodConfigName, t)
	})
}
