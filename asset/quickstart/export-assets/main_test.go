// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := fmt.Sprintf("%s-for-assets", tc.ProjectID)
	env := map[string]string{"GOOGLE_CLOUD_PROJECT": tc.ProjectID}

	ctx := context.Background()

	// Delete the bucket (if it exists) then recreate it.
	testutil.CleanBucket(ctx, t, tc.ProjectID, bucketName)

	m := testutil.BuildMain(t)
	defer m.Cleanup()

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		out, _, err := m.Run(env, 60*time.Second)
		if err != nil {
			r.Errorf("failed to run: %v", err)
			return
		}

		got := string(out)
		want := fmt.Sprintf(`"gs://%s/my-assets.txt`, bucketName)
		if !strings.Contains(got, want) {
			r.Errorf("stdout returned %s, wanted to contain %s", got, want)
		}
	})

}
