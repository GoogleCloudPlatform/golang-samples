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

func TestCreateAkamaiCDNKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	// Create a unique name for the CDN key that incorporates a timestamp.
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}

	// Create a random token key for the CDN key. It is not validated.
	akamaiTokenKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	akamaiCDNKeyID := fmt.Sprintf("%s-%s", akamaiCDNKeyID, uuid)
	akamaiCDNKeyName := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", tc.ProjectID, location, akamaiCDNKeyID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := createCDNKeyAkamai(&buf, tc.ProjectID, akamaiCDNKeyID, hostname, akamaiTokenKey); err != nil {
			r.Errorf("createCDNKeyAkamai got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, akamaiCDNKeyName) {
			r.Errorf("createCDNKeyAkamai got: %v Want to contain: %v", got, akamaiCDNKeyName)
		}
	})
	teardownTestCreateAkamaiCDNKey(akamaiCDNKeyName, t)
}

func teardownTestCreateAkamaiCDNKey(testCDNKeyName string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.DeleteCdnKeyRequest{
		Name: testCDNKeyName,
	}
	op, err := client.DeleteCdnKey(ctx, req)
	if err != nil {
		t.Errorf("client.DeleteCdnKey: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		t.Error(err)
	}
}
