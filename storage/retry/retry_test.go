// Copyright 2021 Google LLC
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

package retry

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestConfigureRetries(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, "storage-buckets-test")
	object := "foo.txt"

	// Upload test object to delete in sample.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, strings.NewReader("hello world")); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}

	var buf bytes.Buffer
	if err := configureRetries(&buf, bucket, object); err != nil {
		t.Errorf("configureRetries: %v", err)
	}

	if got, want := buf.String(), "deleted with a customized retry"; !strings.Contains(got, want) {
		t.Errorf("configureRetries: got %q; want to contain %q", got, want)
	}

}
