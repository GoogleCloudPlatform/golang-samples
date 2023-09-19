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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"strings"
	"testing"
	"time"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

func setupTestGetSlate(slateID string, t *testing.T) func() {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}

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

	return func() {
		req := &stitcherstreampb.DeleteSlateRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID),
		}
		_, err := client.DeleteSlate(ctx, req)
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

func TestGetSlate(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	slateID := "go-get-test-slate"
	teardown := setupTestGetSlate(slateID, t)
	t.Cleanup(teardown)

	slateName := fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := getSlate(buf, tc.ProjectID, slateID); err != nil {
			r.Errorf("getSlate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, slateName) {
			r.Errorf("getSlate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, slateName)
		}
	})
	buf.Reset()
}
