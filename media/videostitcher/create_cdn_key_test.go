// Copyright 2023 Google LLC
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

func TestCreateMediaCDNKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	// Create a unique name for the CDN key that incorporates a timestamp.
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}

	// Create a random private key for the CDN key. It is not validated.
	mediaCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	mediaCDNKeyID := fmt.Sprintf("%s-%s", mediaCDNKeyIDPrefix, uuid)
	mediaCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, mediaCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := createCDNKey(&buf, tc.ProjectID, mediaCDNKeyID, mediaCDNPrivateKey, true); err != nil {
			r.Errorf("createCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, mediaCDNKeyName) {
			r.Errorf("createCDNKey got: %v Want to contain: %v", got, mediaCDNKeyName)
		}
	})

	t.Cleanup(func() {
		deleteTestCDNKey(mediaCDNKeyName, t)
	})
}

func TestCreateCloudCDNKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	// Create a unique name for the CDN key that incorporates a timestamp.
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}

	// Create a random private key for the CDN key. It is not validated.
	cloudCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	cloudCDNKeyID := fmt.Sprintf("%s-%s", cloudCDNKeyIDPrefix, uuid)
	cloudCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, cloudCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := createCDNKey(&buf, tc.ProjectID, cloudCDNKeyID, cloudCDNPrivateKey, false); err != nil {
			r.Errorf("createCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cloudCDNKeyName) {
			r.Errorf("createCDNKey got: %v Want to contain: %v", got, cloudCDNKeyName)
		}
	})

	t.Cleanup(func() {
		deleteTestCDNKey(cloudCDNKeyName, t)
	})
}
