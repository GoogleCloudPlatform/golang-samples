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

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

func setupTestListCDNKeys(keyID string, t *testing.T) func() {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	// client.Close() is called in the returned function

	// Create a random private key for the CDN key. It is not validated.
	mediaCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateCdnKeyRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		CdnKeyId: keyID,
		CdnKey: &stitcherstreampb.CdnKey{
			CdnKeyConfig: &stitcherstreampb.CdnKey_MediaCdnKey{
				MediaCdnKey: &stitcherstreampb.MediaCdnKey{
					KeyName:    keyName,
					PrivateKey: []byte(mediaCDNPrivateKey),
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

	return func() {
		req := &stitcherstreampb.DeleteCdnKeyRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, keyID),
		}
		_, err := client.DeleteCdnKey(ctx, req)
		if err != nil {
			t.Error(err)
		}
		_, err = op.Wait(ctx)
		if err != nil {
			t.Error(err)
		}
		client.Close()
	}
}

func TestListCDNKeys(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	mediaCDNKeyID := fmt.Sprintf("%s-%s", mediaCDNKeyID, uuid)
	teardown := setupTestListCDNKeys(mediaCDNKeyID, t)
	t.Cleanup(teardown)

	mediaCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, mediaCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := listCDNKeys(&buf, tc.ProjectID); err != nil {
			r.Errorf("listCDNKeys got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, mediaCDNKeyName) {
			r.Errorf("listCDNKeys got: %v Want to contain: %v", got, mediaCDNKeyName)
		}
	})
}
