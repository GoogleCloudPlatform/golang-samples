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
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestBlurOffensiveImages(t *testing.T) {
	tc := testutil.SystemTest(t)

	inputBucket := testutil.CreateTestBucket(
		context.Background(),
		t, storageClient, tc.ProjectID, "test-blur-input")
	outputBucket := testutil.CreateTestBucket(
		context.Background(),
		t, storageClient, tc.ProjectID, "test-blur-output")
	oldEnvValue := os.Getenv("BLURRED_BUCKET_NAME")
	os.Setenv("BLURRED_BUCKET_NAME", outputBucket)
	defer os.Setenv("BLURRED_BUCKET_NAME", oldEnvValue)

	e := &storagedata.StorageObjectData{
		Bucket: inputBucket,
		Name:   "zombie.jpg",
	}
	ctx := context.Background()

	inputBlob := storageClient.Bucket(inputBucket).Object(e.Name)
	if _, err := inputBlob.Attrs(ctx); err != nil {
		// input blob does not exist, so upload it.
		bw := inputBlob.NewWriter(context.Background())
		// TODO(muncus): use os.Readfile when we're on go1.16+
		// Note: Open() error will also surface on ReadAll(), so we only check once.
		f, _ := os.Open("zombie.jpg")
		zbytes, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatalf("could not read input file: %v", err)
		}
		// Write is actually performed in Close(), so we check for errors there.
		bw.Write(zbytes)
		if err := bw.Close(); err != nil {
			t.Fatalf("failed to upload file: %v", err)
		}
	}

	outputBlob := storageClient.Bucket(outputBucket).Object(e.Name)

	jsonData, err := protojson.Marshal(e)
	if err != nil {
		t.Fatalf("protojson.Marshal: %v", err)
	}

	ce := cloudevents.NewEvent()
	ce.SetData("application/json", jsonData)
	err = blurOffensiveImages(ctx, ce)
	defer outputBlob.Delete(ctx)
	if err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error: %v", e, err)
	}

	if _, err := outputBlob.Attrs(ctx); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error when checking output: %v", e, err)
	}
}
