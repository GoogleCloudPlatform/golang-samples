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

func setupTestGetLiveConfig(slateID, liveConfigID string, t *testing.T) func() {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	// client.Close() is called in the returned function

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateSlateRequest{
		Parent:  fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		SlateId: slateID,
		Slate: &stitcherstreampb.Slate{
			Uri: slateURI,
		},
	}
	op, err := client.CreateSlate(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = op.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}

	req2 := &stitcherstreampb.CreateLiveConfigRequest{
		Parent:       fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		LiveConfigId: liveConfigID,
		LiveConfig: &stitcherstreampb.LiveConfig{
			SourceUri:       liveURI,
			AdTagUri:        liveAdTagURI,
			AdTracking:      stitcherstreampb.AdTracking_SERVER,
			StitchingPolicy: stitcherstreampb.LiveConfig_CUT_CURRENT,
			DefaultSlate:    fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID),
		},
	}
	op2, err2 := client.CreateLiveConfig(ctx, req2)
	if err2 != nil {
		t.Fatal(err2)
	}
	_, err2 = op2.Wait(ctx)
	if err2 != nil {
		t.Fatal(err2)
	}

	return func() {
		req := &stitcherstreampb.DeleteLiveConfigRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/liveConfigs/%s", tc.ProjectID, location, liveConfigID),
		}
		op, err := client.DeleteLiveConfig(ctx, req)
		if err != nil {
			t.Errorf("client.DeleteLiveConfig: %v", err)
		}
		err = op.Wait(ctx)
		if err != nil {
			t.Error(err)
		}

		req2 := &stitcherstreampb.DeleteSlateRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID),
		}
		_, err2 := client.DeleteSlate(ctx, req2)
		if err2 != nil {
			t.Error(err2)
		}
		err2 = op.Wait(ctx)
		if err2 != nil {
			t.Error(err2)
		}
		client.Close()
	}
}

func TestGetLiveConfig(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	slateID := fmt.Sprintf("%s-%s", slateID, uuid)
	liveConfigID := fmt.Sprintf("%s-%s", liveConfigID, uuid)
	teardown := setupTestGetLiveConfig(slateID, liveConfigID, t)
	t.Cleanup(teardown)

	liveConfigName := fmt.Sprintf("projects/%s/locations/%s/liveConfigs/%s", tc.ProjectID, location, liveConfigID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := getLiveConfig(&buf, tc.ProjectID, liveConfigID); err != nil {
			r.Errorf("getLiveConfig got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, liveConfigName) {
			r.Errorf("getLiveConfig got: %v Want to contain: %v", got, liveConfigName)
		}
	})
}
