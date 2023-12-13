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
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func setupTestUpdateCloudCDNKey(keyID string, t *testing.T) {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	// Create a random private key for the CDN key. It is not validated.
	cloudCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateCdnKeyRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		CdnKeyId: keyID,
		CdnKey: &stitcherstreampb.CdnKey{
			CdnKeyConfig: &stitcherstreampb.CdnKey_GoogleCdnKey{
				GoogleCdnKey: &stitcherstreampb.GoogleCdnKey{
					KeyName:    keyName,
					PrivateKey: []byte(cloudCDNPrivateKey),
				},
			},
			Hostname: hostname,
		},
	}
	op, err := client.CreateCdnKey(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = op.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateCloudCDNKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	updatedCloudCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	cloudCDNKeyID := fmt.Sprintf("%s-%s", cloudCDNKeyIDPrefix, uuid)
	setupTestUpdateCloudCDNKey(cloudCDNKeyID, t)

	cloudCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, cloudCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := updateCDNKey(&buf, tc.ProjectID, cloudCDNKeyID, updatedCloudCDNPrivateKey, false); err != nil {
			r.Errorf("updateCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, cloudCDNKeyName) {
			r.Errorf("updateCDNKey got: %v Want to contain: %v", got, cloudCDNKeyName)
		}
		if got := buf.String(); !strings.Contains(got, updatedHostname) {
			r.Errorf("updateCDNKey got: %v Want to contain: %v", got, updatedHostname)
		}
	})

	t.Cleanup(func() {
		deleteTestCDNKey(cloudCDNKeyName, t)
	})
}

func TestUpdateMediaCDNKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	updatedMediaCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	mediaCDNKeyID := fmt.Sprintf("%s-%s", mediaCDNKeyIDPrefix, uuid)
	createTestMediaCDNKey(mediaCDNKeyID, t)

	mediaCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, mediaCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := updateCDNKey(&buf, tc.ProjectID, mediaCDNKeyID, updatedMediaCDNPrivateKey, true); err != nil {
			r.Errorf("updateCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, mediaCDNKeyName) {
			r.Errorf("updateCDNKey got: %v Want to contain: %v", got, mediaCDNKeyName)
		}
		if got := buf.String(); !strings.Contains(got, updatedHostname) {
			r.Errorf("updateCDNKey got: %v Want to contain: %v", got, updatedHostname)
		}
	})

	t.Cleanup(func() {
		deleteTestCDNKey(mediaCDNKeyName, t)
	})
}
