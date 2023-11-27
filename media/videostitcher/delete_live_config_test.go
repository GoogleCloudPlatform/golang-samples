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

func setupTestDeleteLiveConfig(slateID, liveConfigID string, t *testing.T) {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	// Create a slate.
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

	// Create a live config to delete.
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
}

func TestDeleteLiveConfig(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	slateID := fmt.Sprintf("%s-%s", slateID, uuid)
	slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
	liveConfigID := fmt.Sprintf("%s-%s", liveConfigID, uuid)
	setupTestDeleteLiveConfig(slateID, liveConfigID, t)

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteLiveConfig(&buf, tc.ProjectID, liveConfigID); err != nil {
			r.Errorf("deleteLiveConfig got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteLiveConfigResponse) {
			r.Errorf("deleteLiveConfig got: %v Want to contain: %v", got, deleteLiveConfigResponse)
		}
	})

	t.Cleanup(func() {
		teardownTestDeleteLiveConfig(slateName, t)
	})
}

func teardownTestDeleteLiveConfig(slateName string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	// Delete the slate.
	req := &stitcherstreampb.DeleteSlateRequest{
		Name: slateName,
	}
	op, err := client.DeleteSlate(ctx, req)
	if err != nil {
		t.Errorf("client.DeleteSlate: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		t.Error(err)
	}
}
