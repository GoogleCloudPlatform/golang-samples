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

func TestCreateSlate(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	uuid, err := getUUID()
	if err != nil {
		t.Fatalf("getUUID err: %v", err)
	}
	slateID := fmt.Sprintf("%s-%s", slateID, uuid)
	slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := createSlate(&buf, tc.ProjectID, slateID, slateURI); err != nil {
			r.Errorf("createSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			r.Errorf("createSlate got: %v Want to contain: %v", got, slateName)
		}
	})
	t.Cleanup(func() {
		teardownTestCreateSlate(slateName, t)
	})
}

func teardownTestCreateSlate(slateName string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

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
