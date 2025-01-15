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

func setupTestUpdateCDNKeyAkamai(keyID string, t *testing.T) {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	// Create a random token key for the CDN key. It is not validated.
	akamaiTokenKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateCdnKeyRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		CdnKeyId: keyID,
		CdnKey: &stitcherstreampb.CdnKey{
			CdnKeyConfig: &stitcherstreampb.CdnKey_AkamaiCdnKey{
				AkamaiCdnKey: &stitcherstreampb.AkamaiCdnKey{
					TokenKey: []byte(akamaiTokenKey),
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

func TestUpdateCDNKeyAkamai(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	akamaiTokenKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	akamaiCDNKeyID := fmt.Sprintf("%s-%s", akamaiCDNKeyIDPrefix, uuid)
	setupTestUpdateCDNKeyAkamai(akamaiCDNKeyID, t)

	akamaiCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, akamaiCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := updateCDNKeyAkamai(&buf, tc.ProjectID, akamaiCDNKeyID, akamaiTokenKey); err != nil {
			r.Errorf("updateCDNKey got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, akamaiCDNKeyName) {
			r.Errorf("updateCDNKey got: %v Want to contain: %v", got, akamaiCDNKeyName)
		}
		if got := buf.String(); !strings.Contains(got, updatedHostname) {
			r.Errorf("updateCDNKey got: %v Want to contain: %v", got, updatedHostname)
		}
	})

	t.Cleanup(func() {
		deleteTestCDNKey(akamaiCDNKeyName, t)
	})
}
