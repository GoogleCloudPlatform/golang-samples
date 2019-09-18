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
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestBlurOffensiveImages(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID not set")
	}

	log.SetOutput(ioutil.Discard)

	outputBucket := projectID + "-test-blurred"
	os.Setenv("BLURRED_BUCKET_NAME", outputBucket)

	e := GCSEvent{
		Bucket: projectID,
		Name:   "functions/zombie.jpg",
	}
	ctx := context.Background()

	inputBlob := storageClient.Bucket(projectID).Object(e.Name)
	if _, err := inputBlob.Attrs(ctx); err != nil {
		t.Skipf("could not get input file: %s: %v", inputBlob.ObjectName(), err)
	}

	b := storageClient.Bucket(outputBucket)
	b.Create(ctx, projectID, nil)
	outputBlob := b.Object(e.Name)
	outputBlob.Delete(ctx) // Ensure the output file doesn't already exist.

	if err := BlurOffensiveImages(ctx, e); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error: %v", e, err)
	}

	if _, err := outputBlob.Attrs(ctx); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error when checking output: %v", e, err)
	}
	outputBlob.Delete(ctx)
}
