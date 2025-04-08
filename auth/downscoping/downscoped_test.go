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

package downscopedoverview

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"golang.org/x/oauth2/google/downscope"
)

func TestInitializeCredentials(t *testing.T) {
	accessBoundary := []downscope.AccessBoundaryRule{
		{
			AvailableResource:    "//storage.googleapis.com/projects/_/buckets/foo",
			AvailablePermissions: []string{"inRole:roles/storage.objectViewer"},
		},
	}
	err := initializeCredentials(accessBoundary)
	if err != nil {
		t.Errorf("got %v; wanted nil", err)
	}
}

func TestReadObjectContents(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	randSuffix, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("Failed to generate random UUID suffix: %v", err)
	}
	bucketName := "bucket-downscoping-test-golang-" + randSuffix.String()[:8]
	objectName := "object-downscoping-test-golang-" + randSuffix.String()[:8]
	content := "CONTENT"
	bucket := client.Bucket(bucketName)

	if err := bucket.Create(ctx, tc.ProjectID, nil); err != nil {
		t.Errorf("Failed to create bucket: %v", err)
		return
	}
	defer bucket.Delete(ctx)
	obj := bucket.Object(objectName)
	write := obj.NewWriter(ctx)
	_, err = io.Copy(write, strings.NewReader(content))
	if err != nil {
		t.Errorf("failed to copy the text: %v", err)
		return
	}
	defer obj.Delete(ctx)

	err = write.Close()
	if err != nil {
		t.Errorf("failed to close the writer: %v", err)
		return
	}

	buf := new(bytes.Buffer)
	err = getObjectContents(buf, bucketName, objectName)
	if err != nil {
		t.Errorf("failed to retrieve object contents: %v", err)
		return
	}
	out := buf.String()
	if out != content {
		t.Errorf("got %v but want "+content, out)
		return
	}
}
