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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func createBucket(t *testing.T, projectID, bucketName string) error {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("storage.NewClient: %v", err)
	}

	bucket := storageClient.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		t.Errorf("Bucket(%q).Create: %v", bucketName, err)
	}

	return nil
}

func deleteBucket(t *testing.T, bucketName string) error {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("storage.NewClient: %v", err)
	}

	bucket := storageClient.Bucket(bucketName)
	if err := bucket.Delete(ctx); err != nil {
		t.Errorf("Bucket(%q).Delete: %v", bucketName, err)
	}

	return nil
}

func TestUsageExportSnippets(t *testing.T) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	bucketName := "test-bucket-" + fmt.Sprint(seededRand.Int())

	createBucket(t, tc.ProjectID, bucketName)
	defer deleteBucket(t, bucketName)

	buf := &bytes.Buffer{}

	if err := setUsageExportBucket(buf, tc.ProjectID, bucketName, ""); err != nil {
		t.Errorf("setUsageExportBucket got err: %v", err)
	}

	expectedResult := "Setting reportNamePrefix to empty value causes the report to have the default prefix value `usage_gce`"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("setUsageExportBucket got %q, want %q", got, expectedResult)
	}
	expectedResult = "Usage export bucket has been set"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("setUsageExportBucket got %q, want %q", got, expectedResult)
	}

	projectsClient, err := compute.NewProjectsRESTClient(ctx)
	if err != nil {
		t.Errorf("NewProjectsRESTClient: %v", err)
	}
	defer projectsClient.Close()

	req := &computepb.GetProjectRequest{
		Project: tc.ProjectID,
	}

	project, err := projectsClient.Get(ctx, req)
	if err != nil {
		t.Errorf("Project get request: %v", err)
	}

	usageExportLocation := project.GetUsageExportLocation()

	if usageExportLocation.GetBucketName() != bucketName {
		t.Errorf("Got: %s; want %s", *usageExportLocation.BucketName, bucketName)
	}
	if usageExportLocation.GetReportNamePrefix() != "" {
		t.Errorf("Got: %s; want %q", *usageExportLocation.BucketName, "")
	}

	buf.Reset()

	if err := getUsageExportBucket(buf, tc.ProjectID); err != nil {
		t.Errorf("setUsageExportBucket got err: %v", err)
	}

	expectedResult = "Report name prefix not set, replacing with default value of `usage_gce`"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getUsageExportBucket got %q, want %q", got, expectedResult)
	}

	expectedResult = "Returned ReportNamePrefix: usage_gce"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getUsageExportBucket got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := disableUsageExport(buf, tc.ProjectID); err != nil {
		t.Errorf("disableUsageExport got err: %v", err)
	}

	expectedResult = "Usage export bucket has been set"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getUsageExportBucket got %q, want %q", got, expectedResult)
	}

	req = &computepb.GetProjectRequest{
		Project: tc.ProjectID,
	}

	project, err = projectsClient.Get(ctx, req)
	if err != nil {
		t.Errorf("Project get request: %v", err)
	}

	if project.GetUsageExportLocation() != nil {
		t.Errorf("UsageExportLocation should be nil")
	}
}
