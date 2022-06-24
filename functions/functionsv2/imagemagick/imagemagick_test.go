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

package imagemagick

import (
	"context"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestBlurOffensiveImages(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID not set")
	}

	outputBucket, err := testutil.CreateTestBucket(
		context.Background(),
		t, storageClient, projectID, "test-blur-output")
	if err != nil {
		t.Errorf("failed to create output bucket: %v", err)
	}
	oldEnvValue := os.Getenv("BLURRED_BUCKET_NAME")
	os.Setenv("BLURRED_BUCKET_NAME", outputBucket)
	defer os.Setenv("BLURRED_BUCKET_NAME", oldEnvValue)

	e := GCSEvent{
		Bucket: projectID,
		Name:   "functions/zombie.jpg",
	}
	ctx := context.Background()

	inputBlob := storageClient.Bucket(projectID).Object(e.Name)
	if _, err := inputBlob.Attrs(ctx); err != nil {
		t.Skipf("could not get input file: %s: %v", inputBlob.ObjectName(), err)
	}

	outputBlob := storageClient.Bucket(outputBucket).Object(e.Name)

	ce := cloudevents.NewEvent()
	ce.SetData("application/json", e)
	err = blurOffensiveImages(ctx, ce)
	defer outputBlob.Delete(ctx)
	if err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error: %v", e, err)
	}

	if _, err := outputBlob.Attrs(ctx); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error when checking output: %v", e, err)
	}
}
