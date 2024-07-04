// Copyright 2023 Google LLC
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

package inspect

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectGCSFileSendToScc(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	ctx := context.Background()
	sc, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer sc.Close()

	// Creates a bucket using a function available in testutil.
	bucketNameForInspectGCSSendToScc := testutil.CreateTestBucket(ctx, t, sc, tc.ProjectID, "dlp-test-inspect-prefix")

	// Uploads a file on created bucket.
	filePathtoGCS(t, tc.ProjectID, bucketNameForInspectGCSSendToScc, dirPathForInspectGCSSendToScc)

	gcsPath := fmt.Sprint("gs://" + bucketNameForInspectGCSSendToScc + "/" + dirPathForInspectGCSSendToScc + "/test.txt")

	if err := inspectGCSFileSendToScc(&buf, tc.ProjectID, gcsPath); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("TestInspectGCSFileSendToScc got %q, want %q", got, want)
	}

	// Delete a bucket that has just been created.
	err = testutil.DeleteBucketIfExists(ctx, sc, bucketNameForInspectGCSSendToScc)
	if err != nil {
		t.Fatal(err)
	}
}
