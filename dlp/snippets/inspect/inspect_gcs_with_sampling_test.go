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
	"log"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectGcsFileWithSampling(t *testing.T) {
	tc := testutil.SystemTest(t)
	topicID := "go-lang-dlp-test-bigquery-with-sampling-topic"
	subscriptionID := "go-lang-dlp-test-bigquery-with-sampling-subscription"
	ctx := context.Background()
	sc, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer sc.Close()

	bucketnameForInspectGCSFileWithSampling := testutil.CreateTestBucket(ctx, t, sc, tc.ProjectID, "dlp-test-inspect-prefix")
	GCSUri := "gs://" + bucketnameForInspectGCSFileWithSampling + "/"

	var buf bytes.Buffer
	if err := inspectGcsFileWithSampling(&buf, tc.ProjectID, GCSUri, topicID, subscriptionID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Job Created"; !strings.Contains(got, want) {
		t.Errorf("inspectGcsFileWithSampling got %q, want %q", got, want)
	}
	err = testutil.DeleteBucketIfExists(ctx, sc, bucketnameForInspectGCSFileWithSampling)
	if err != nil {
		t.Fatal(err)
	}

}
